import { Device } from './DeviceList'
import React, { useState } from 'react'
import { Button, FormControl, Grid, InputLabel, MenuItem, Select, TextField } from '@material-ui/core'
import { makeStyles } from '@material-ui/core/styles'
import { SelectInputProps } from '@material-ui/core/Select/SelectInput'
import { Autocomplete } from '@material-ui/lab'
import SubPanel from './SubPanel'

const useStyles = makeStyles((theme) => ({
  settingsPanelRow: {
    marginBottom: theme.spacing(1)
  }
}))


export default function DeviceSettingsPanel(props: { device: Device }) {
  const classes = useStyles()

  const [dev, setDevice] = useState(props.device)
  const [pollError, setPollError] = useState(false)
  const [instanceError, setInstanceError] = useState(false)

  const onSave = () => {
    console.log('Clicked!', dev)
  }

  return <SubPanel heading={'Settings'}>
    <Grid item container spacing={3} className={classes.settingsPanelRow}>
      <Grid item xs={4} sm={3} md={4} lg={3}>
        <InstanceTextField instance={dev.instance} onInstanceChange={(instance, err) => {
          setInstanceError(err)
          setDevice({ ...dev, instance: instance })
        }}/>
      </Grid>
      <Grid item xs={4} sm={3} md={4} lg={3}>
        <TxPowerSelect txPower={dev.txPower} onTxPowerSelected={(e) => {
          setDevice({ ...dev, txPower: parseInt(e.target.value as string) })
        }}/>
      </Grid>
      <Grid item xs={4} sm={3} md={4} lg={3}>
        <PollPeriodAutoComplete pollPeriod={dev.pollPeriod} onPollPeriodChange={(period, err) => {
          setPollError(err)
          if (!err) {
            setDevice({ ...dev, pollPeriod: period })
          }
        }}/>
      </Grid>
    </Grid>
    <Grid item xs={12} className={classes.settingsPanelRow}>
      <Button variant={'outlined'} color={'primary'} disabled={pollError || instanceError} onClick={onSave}>
        Save
      </Button>
    </Grid>
  </SubPanel>
}

function InstanceTextField(props: { instance: string, onInstanceChange: (instance: string, err: boolean) => void }) {
  const regex = /^[\w]{2,4}$/
  const [err, setErr] = useState(false)

  return <TextField label="Instance" error={err} value={props.instance} onChange={(e) => {
    const v = e.target.value
    const hasError = !regex.test(v)
    setErr(hasError)
    props.onInstanceChange(v, hasError)
  }}/>
}

function TxPowerSelect(props: { txPower: number, onTxPowerSelected: SelectInputProps['onChange'] }) {
  return <FormControl fullWidth>
    <InputLabel id="txpower-label">TX Power</InputLabel>
    <Select labelId="txpower-label" defaultValue={props.txPower} onChange={props.onTxPowerSelected}>
      <MenuItem value={8}>8 dBm</MenuItem>
      <MenuItem value={4}>4 dBm</MenuItem>
      <MenuItem value={0}>0 dBm</MenuItem>
      <MenuItem value={-4}>-4 dBm</MenuItem>
      <MenuItem value={-8}>-8 dBm</MenuItem>
      <MenuItem value={-12}>-12 dBm</MenuItem>
      <MenuItem value={-16}>-16 dBm</MenuItem>
      <MenuItem value={-20}>-20 dBm</MenuItem>
    </Select>
  </FormControl>
}

function PollPeriodAutoComplete(props: { pollPeriod: number, onPollPeriodChange: (period: number, err: boolean) => void }) {
  const [err, setErr] = useState(false)
  const regex = /^([\d]+)(ms$| ms$|[\s]*$)/

  return <Autocomplete
    freeSolo
    options={['200 ms', '500 ms', '1000 ms', '2000 ms', '5000 ms', '10000 ms', '15000 ms']}
    value={props.pollPeriod.toString() + ' ms'}
    onInputChange={(e, val, reason) => {
      const matches = val.match(regex)
      if (matches !== null) {
        props.onPollPeriodChange(parseInt(matches[1]), false)
      } else {
        props.onPollPeriodChange(NaN, true)
      }
      setErr(matches === null)
    }}
    filterOptions={(opts, state) => opts}
    renderInput={(params) => (
      <TextField {...params} error={err} label="Poll Period"/>
    )}
  />
}


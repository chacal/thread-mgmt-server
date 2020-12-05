import React, { useState } from 'react'
import { makeStyles } from '@material-ui/core/styles'
import { SelectInputProps } from '@material-ui/core/Select/SelectInput'
import Autocomplete from '@material-ui/lab/Autocomplete'
import SubPanel from './SubPanel'
import AsyncOperationButton from './AsyncOperationButton'
import Grid from '@material-ui/core/Grid'
import TextField from '@material-ui/core/TextField'
import FormControl from '@material-ui/core/FormControl'
import InputLabel from '@material-ui/core/InputLabel'
import Select from '@material-ui/core/Select'
import MenuItem from '@material-ui/core/MenuItem'
import isEqual from 'lodash/isEqual'
import ErrorMessage from './ErrorMessage'

const useStyles = makeStyles((theme) => ({
  settingsPanelRow: {
    marginBottom: theme.spacing(1)
  }
}))

export interface DeviceSettings {
  instance?: string
  txPower?: number
  pollPeriod?: number
}

export default function DeviceSettingsPanel(props: { settings: DeviceSettings, onSaveSettings: (s: DeviceSettings) => Promise<void> }) {
  const classes = useStyles()

  const [settings, setSettings] = useState(props.settings)
  const [pollError, setPollError] = useState(false)
  const [instanceError, setInstanceError] = useState(false)
  const [saveError, setSaveError] = useState('')

  const onClickSave = () => {
    setSaveError('')
    return props.onSaveSettings(settings)
      .catch(err => setSaveError(err.toString()))
  }

  return <SubPanel heading={'Settings'}>
    <Grid item container spacing={3} className={classes.settingsPanelRow}>
      <Grid item xs={4} sm={3} md={4} lg={3}>
        <InstanceTextField instance={settings.instance} onInstanceChange={(instance, err) => {
          setInstanceError(err)
          setSettings({ ...settings, instance: instance })
        }}/>
      </Grid>
      <Grid item xs={4} sm={3} md={4} lg={3}>
        <TxPowerSelect txPower={settings.txPower} onTxPowerSelected={(e) => {
          setSettings({ ...settings, txPower: parseInt(e.target.value as string) })
        }}/>
      </Grid>
      <Grid item xs={4} sm={3} md={4} lg={3}>
        <PollPeriodAutoComplete pollPeriod={settings.pollPeriod} onPollPeriodChange={(period, err) => {
          setPollError(err)
          if (!err) {
            setSettings({ ...settings, pollPeriod: period })
          }
        }}/>
      </Grid>
    </Grid>
    <Grid item xs={12} className={classes.settingsPanelRow}>
      <AsyncOperationButton disabled={pollError || instanceError || isEqual(settings, props.settings)}
                            onClick={onClickSave}>
        Save
      </AsyncOperationButton>
      <ErrorMessage msg={saveError}/>
    </Grid>
  </SubPanel>
}

function InstanceTextField(props: { instance?: string, onInstanceChange: (instance: string, err: boolean) => void }) {
  const regex = /^[\w]{2,4}$/
  const [err, setErr] = useState(false)

  return <TextField label="Instance" error={err} value={props.instance} onChange={(e) => {
    const v = e.target.value
    const hasError = !regex.test(v)
    setErr(hasError)
    props.onInstanceChange(v, hasError)
  }}/>
}

function TxPowerSelect(props: { txPower?: number, onTxPowerSelected: SelectInputProps['onChange'] }) {
  return <FormControl fullWidth>
    <InputLabel id="txpower-label">TX Power</InputLabel>
    <Select labelId="txpower-label" defaultValue={props.txPower ? props.txPower : ''}
            onChange={props.onTxPowerSelected}>
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

function PollPeriodAutoComplete(props: { pollPeriod?: number, onPollPeriodChange: (period: number, err: boolean) => void }) {
  const [err, setErr] = useState(false)
  const regex = /^([\d]+)(ms$| ms$|[\s]*$)/

  return <Autocomplete
    freeSolo
    options={['200 ms', '500 ms', '1000 ms', '2000 ms', '5000 ms', '10000 ms', '15000 ms']}
    value={props.pollPeriod ? props.pollPeriod.toString() + ' ms' : ''}
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


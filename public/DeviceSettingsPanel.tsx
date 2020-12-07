import React, { ChangeEvent, useState } from 'react'
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
import { DeviceDefaults } from './DeviceList'
import FormHelperText from '@material-ui/core/FormHelperText'

const useStyles = makeStyles((theme) => ({
  settingsPanelRow: {
    marginBottom: theme.spacing(1)
  }
}))

export default function DeviceSettingsPanel(props: { defaults: DeviceDefaults, onSaveDefaults: (s: DeviceDefaults) => Promise<void> }) {
  const classes = useStyles()

  const [defaults, setDefaults] = useState(props.defaults)
  const [pollError, setPollError] = useState(false)
  const [instanceError, setInstanceError] = useState(false)
  const [saveError, setSaveError] = useState('')

  const onInstanceChange = (instance: string, err: boolean) => {
    setInstanceError(err)
    if (instance !== '') {
      setDefaults({ ...defaults, instance })
    } else {
      setDefaults(prev => {
        const { instance, ...rest } = prev
        return rest
      })
    }
  }

  const onTxPowerSelected = (e: ChangeEvent<HTMLSelectElement>) => {
    setDefaults({ ...defaults, txPower: parseInt(e.target.value) })
  }

  const onPollPeriodChange = (period: number | undefined, err: boolean) => {
    setPollError(err)
    if (!err) {
      if (period !== undefined) {
        setDefaults({ ...defaults, pollPeriod: period })
      } else {
        setDefaults(prev => {
          const { pollPeriod, ...rest } = prev
          return rest
        })
      }
    }
  }

  const onClickSave = () => {
    setSaveError('')
    return props.onSaveDefaults(defaults)
      .catch(err => setSaveError(err.toString()))
  }

  const isSaveDisabled = () => pollError || instanceError || isEqual(defaults, props.defaults)

  return <SubPanel heading={'Settings'}>
    <Grid item container spacing={3} className={classes.settingsPanelRow}>
      <Grid item xs={4} sm={3} md={4} lg={3}>
        <InstanceTextField instance={defaults.instance} onInstanceChange={onInstanceChange}/>
      </Grid>
      <Grid item xs={4} sm={3} md={4} lg={3}>
        <TxPowerSelect txPower={defaults.txPower} onTxPowerSelected={onTxPowerSelected}/>
      </Grid>
      <Grid item xs={4} sm={3} md={4} lg={3}>
        <PollPeriodAutoComplete pollPeriod={defaults.pollPeriod} onPollPeriodChange={onPollPeriodChange}/>
      </Grid>
    </Grid>
    <Grid item xs={12} className={classes.settingsPanelRow}>
      <AsyncOperationButton disabled={isSaveDisabled()} onClick={onClickSave}>Save</AsyncOperationButton>
      <ErrorMessage msg={saveError}/>
    </Grid>
  </SubPanel>
}

function InstanceTextField(props: { instance?: string, onInstanceChange: (instance: string, err: boolean) => void }) {
  const regex = /^[\w]{2,4}$|^$/
  const [err, setErr] = useState(false)

  const onChange = (e: ChangeEvent<HTMLInputElement>) => {
    const v = e.target.value
    const hasError = !regex.test(v)
    setErr(hasError)
    props.onInstanceChange(v, hasError)
  }

  return <TextField label="Instance" error={err} value={props.instance ? props.instance : ''} onChange={onChange}/>
}

function TxPowerSelect(props: { txPower?: number, onTxPowerSelected: SelectInputProps['onChange'] }) {
  return <FormControl fullWidth>
    <InputLabel id="txpower-label">TX Power</InputLabel>
    <Select labelId="txpower-label" value={props.txPower !== undefined ? props.txPower : ''}
            onChange={props.onTxPowerSelected}>
      <MenuItem value={8}>8</MenuItem>
      <MenuItem value={4}>4</MenuItem>
      <MenuItem value={0}>0</MenuItem>
      <MenuItem value={-4}>-4</MenuItem>
      <MenuItem value={-8}>-8</MenuItem>
      <MenuItem value={-12}>-12</MenuItem>
      <MenuItem value={-16}>-16</MenuItem>
      <MenuItem value={-20}>-20</MenuItem>
    </Select>
    <FormHelperText>dBm</FormHelperText>
  </FormControl>
}

function PollPeriodAutoComplete(props: { pollPeriod?: number, onPollPeriodChange: (period: number | undefined, err: boolean) => void }) {
  const [err, setErr] = useState(false)

  const onInputChange = (e: React.ChangeEvent, val: string) => {
    const valid = isValidPollPeriod(val)
    if (valid) {
      props.onPollPeriodChange(parseValidPollPeriod(val), false)
    } else {
      props.onPollPeriodChange(undefined, true)
    }
    setErr(!valid)
  }

  return <Autocomplete
    freeSolo
    options={['200', '500', '1000', '2000', '5000', '10000', '15000']}
    value={props.pollPeriod !== undefined ? props.pollPeriod.toString() : ''}
    onInputChange={onInputChange}
    filterOptions={(opts, state) => opts}
    renderInput={(params) => (
      <TextField {...params} error={err} label="Poll Period" helperText="ms"/>
    )}
  />
}

const pollPeriodRegex = /^([\d]+)([\s]*$)/

function isValidPollPeriod(val: string): boolean {
  const matches = val.match(pollPeriodRegex)
  return (matches !== null && parseInt(matches[1]) > 0) || val === ''
}

function parseValidPollPeriod(val: string): number | undefined {
  const matches = val.match(pollPeriodRegex)
  if (matches !== null) {
    return parseInt(matches[1])
  } else {
    return undefined
  }
}
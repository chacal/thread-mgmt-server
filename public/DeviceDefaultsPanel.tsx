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
import StatusMessage, { EmptyStatus, Status } from './StatusMessage'
import { DeviceDefaults } from './DeviceList'
import InputAdornment from '@material-ui/core/InputAdornment'
import { postJSON } from './DeviceListItem'

const useStyles = makeStyles((theme) => ({
  defaultsPanelInputs: {
    marginBottom: theme.spacing(1)
  },
  txPowerSelectAdornment: {
    position: 'absolute',
    padding: 0,
    right: '16px',
    top: 'calc(50%)',
  }
}))

export default function DeviceDefaultsPanel(props: { defaults: DeviceDefaults, deviceId: string, mainIp: string | undefined, onSaveDefaults: (s: DeviceDefaults) => Promise<void> }) {
  const classes = useStyles()

  const [defaults, setDefaults] = useState(props.defaults)
  const [pollError, setPollError] = useState(false)
  const [instanceError, setInstanceError] = useState(false)
  const [status, setStatus] = useState(EmptyStatus)

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

  const setErrorStatus = (err: any) => setStatus({ msg: err.toString(), isError: true, showProgress: false })

  const onClickSave = () => {
    setStatus(EmptyStatus)
    return props.onSaveDefaults(defaults)
      .catch(setErrorStatus)
  }

  const onClickPush = () => {
    setStatus({ msg: 'Pushing defaults..', isError: false, showProgress: true })
    return postJSON('/v1/devices/' + props.deviceId + '/push', { address: props.mainIp })
      .then(() => setStatus(EmptyStatus))
      .catch(setErrorStatus)
  }

  const isSaveDisabled = () => pollError || instanceError || isEqual(defaults, props.defaults)
  const isPushDisabled = () => props.mainIp === undefined || !isEqual(defaults, props.defaults)

  return <SubPanel heading={'Defaults'}>
    <Grid item container spacing={3} className={classes.defaultsPanelInputs}>
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
    <Grid item container spacing={2} xs={12} alignItems={'center'}>
      <Grid item>
        <AsyncOperationButton disabled={isSaveDisabled()} onClick={onClickSave}>Save</AsyncOperationButton>
      </Grid>
      <Grid item>
        <AsyncOperationButton disabled={isPushDisabled()} onClick={onClickPush}>Push</AsyncOperationButton>
      </Grid>
    </Grid>
    <Grid item xs={12}>
      <StatusMessage {...status}/>
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

  return <TextField
    label="Instance"
    error={err}
    value={props.instance ? props.instance : ''}
    onChange={onChange}
    InputLabelProps={{ shrink: true }}
  />
}

function TxPowerSelect(props: { txPower?: number, onTxPowerSelected: SelectInputProps['onChange'] }) {
  const classes = useStyles()

  return <FormControl fullWidth>
    <InputLabel shrink={true} id="txpower-label">TX Power</InputLabel>
    <Select labelId="txpower-label"
            value={props.txPower !== undefined ? props.txPower : ''}
            onChange={props.onTxPowerSelected}
            endAdornment={
              <InputAdornment position="start" className={classes.txPowerSelectAdornment}>dBm</InputAdornment>
            }
    >
      <MenuItem value={8}>8</MenuItem>
      <MenuItem value={4}>4</MenuItem>
      <MenuItem value={0}>0</MenuItem>
      <MenuItem value={-4}>-4</MenuItem>
      <MenuItem value={-8}>-8</MenuItem>
      <MenuItem value={-12}>-12</MenuItem>
      <MenuItem value={-16}>-16</MenuItem>
      <MenuItem value={-20}>-20</MenuItem>
    </Select>
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
      <TextField
        {...params}
        error={err}
        label="Poll Period"
        InputProps={{ endAdornment: <InputAdornment position="end">ms</InputAdornment> }}
        InputLabelProps={{ shrink: true }}
      />
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
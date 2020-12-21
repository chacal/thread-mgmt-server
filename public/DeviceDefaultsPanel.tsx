import React, { ChangeEvent, useState } from 'react'
import { makeStyles } from '@material-ui/core/styles'
import { SelectInputProps } from '@material-ui/core/Select/SelectInput'
import SubPanel from './SubPanel'
import AsyncOperationButton from './AsyncOperationButton'
import Grid from '@material-ui/core/Grid'
import TextField from '@material-ui/core/TextField'
import FormControl from '@material-ui/core/FormControl'
import InputLabel from '@material-ui/core/InputLabel'
import Select from '@material-ui/core/Select'
import MenuItem from '@material-ui/core/MenuItem'
import isEqual from 'lodash/isEqual'
import StatusMessage, { EmptyStatus } from './StatusMessage'
import { DeviceDefaults } from './DeviceList'
import InputAdornment from '@material-ui/core/InputAdornment'
import { postJSON } from './DeviceListItem'

const useStyles = makeStyles((theme) => ({
  defaultsPanelInputs: {
    marginBottom: theme.spacing(1)
  },
  selectEndAdornment: {
    position: 'absolute',
    padding: 0,
    right: '16px',
    top: 'calc(50%)',
  }
}))

export default function DeviceDefaultsPanel(props: { defaults: DeviceDefaults, deviceId: string, mainIp: string, onSaveDefaults: (s: DeviceDefaults) => Promise<void> }) {
  const classes = useStyles()

  const [defaults, setDefaults] = useState(props.defaults)
  const [status, setStatus] = useState(EmptyStatus)

  const onInstanceChange = (e: ChangeEvent<HTMLInputElement>) => {
    setDefaults({ ...defaults, instance: e.target.value })
  }

  const onTxPowerSelected = (e: ChangeEvent<HTMLSelectElement>) => {
    setDefaults({ ...defaults, txPower: parseInt(e.target.value) })
  }

  const onPollPeriodChange = (e: ChangeEvent<HTMLSelectElement>) => {
    setDefaults({ ...defaults, pollPeriod: parseInt(e.target.value) })
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

  const isSaveDisabled = () => !isValidInstance(defaults.instance) || isEqual(defaults, props.defaults)
  const isPushDisabled = () => props.mainIp === '' || !isEqual(defaults, props.defaults)

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

function InstanceTextField(props: { instance: string, onInstanceChange: (e: ChangeEvent<HTMLInputElement>) => void }) {
  return <TextField label="Instance" error={!isValidInstance(props.instance)} value={props.instance}
                    onChange={props.onInstanceChange} InputLabelProps={{ shrink: true }}/>
}

function TxPowerSelect(props: { txPower: number, onTxPowerSelected: SelectInputProps['onChange'] }) {
  const classes = useStyles()

  return <FormControl fullWidth>
    <InputLabel shrink={true} id="txpower-label">TX Power</InputLabel>
    <Select labelId="txpower-label"
            value={props.txPower}
            onChange={props.onTxPowerSelected}
            endAdornment={
              <InputAdornment position="start" className={classes.selectEndAdornment}>dBm</InputAdornment>
            }
    >
      {[8, 4, 0, -4, -8, -12, -16, -20].map(v => <MenuItem key={v} value={v}>{v}</MenuItem>)}
    </Select>
  </FormControl>
}

function PollPeriodAutoComplete(props: { pollPeriod: number, onPollPeriodChange: SelectInputProps['onChange'] }) {
  const classes = useStyles()

  return <FormControl fullWidth>
    <InputLabel shrink={true} id="poll-period-label">Poll Period</InputLabel>
    <Select labelId="poll-period-label"
            value={props.pollPeriod}
            onChange={props.onPollPeriodChange}
            endAdornment={
              <InputAdornment position="start" className={classes.selectEndAdornment}>ms</InputAdornment>
            }
    >
      {[50, 100, 200, 500, 1000, 2000, 5000, 10000, 15000].map(v => <MenuItem key={v} value={v}>{v}</MenuItem>)}
    </Select>
  </FormControl>
}

function isValidInstance(instance: string) {
  const regex = /^[\w]{2,4}$/
  return regex.test(instance)
}
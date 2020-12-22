import { DeviceConfig } from './DeviceList'
import React, { ChangeEvent, useState } from 'react'
import SubPanel from './SubPanel'
import InputLabel from '@material-ui/core/InputLabel'
import Select from '@material-ui/core/Select'
import MenuItem from '@material-ui/core/MenuItem'
import FormControl from '@material-ui/core/FormControl'
import Grid from '@material-ui/core/Grid'
import AsyncOperationButton from './AsyncOperationButton'
import StatusMessage, { EmptyStatus } from './StatusMessage'
import isEqual from 'lodash/isEqual'
import { makeStyles } from '@material-ui/core/styles'
import Checkbox from '@material-ui/core/Checkbox'
import FormControlLabel from '@material-ui/core/FormControlLabel'
import InputAdornment from '@material-ui/core/InputAdornment'
import Typography from '@material-ui/core/Typography'

const useStyles = makeStyles((theme) => ({
  configPanelRow: {
    marginBottom: theme.spacing(1)
  },
  pollIntervalSelectAdornment: {
    position: 'absolute',
    padding: 0,
    right: '16px',
    top: 'calc(50%)',
  },
  checkBoxRoot: {
    padding: theme.spacing(0.5, 1)
  }
}))

export default function DeviceConfigPanel(props: { config: DeviceConfig, addresses: string[] | undefined, onSaveConfig: (c: DeviceConfig) => Promise<void> }) {
  const classes = useStyles()
  const addresses = props.addresses !== undefined ? props.addresses : []

  const [config, setConfig] = useState(props.config)
  const [status, setStatus] = useState(EmptyStatus)

  const onMainIPSelected = (mainIp: string) => setConfig({ ...config, mainIp })

  const onPollConfigChanged = (enabled: boolean, interval: number) => {
    setConfig({ ...config, statePollingEnabled: enabled, statePollingIntervalSec: interval })
  }

  const onClickSave = () => {
    setStatus(EmptyStatus)
    return props.onSaveConfig(config)
      .catch(err => setStatus({ msg: err.toString(), isError: true, showProgress: false }))
  }

  const isSaveDisabled = () => isEqual(config, props.config)

  return <SubPanel heading={'Config'}>
    <Grid item container spacing={3} className={classes.configPanelRow}>
      <Grid item xs={12} sm={9} md={10} lg={8}>
        <MainIPSelect mainIp={config.mainIp} addresses={addresses} onMainIPSelected={onMainIPSelected}/>
      </Grid>
    </Grid>
    <Grid item container spacing={3} className={classes.configPanelRow}>
      <StatePollingControls
        disabled={config.mainIp === ''}
        pollingEnabled={config.statePollingEnabled}
        pollingIntervalSec={config.statePollingIntervalSec}
        onPollChange={onPollConfigChanged}
      />
    </Grid>
    <Grid item container spacing={2} xs={12}>
      <Grid item>
        <AsyncOperationButton disabled={isSaveDisabled()} onClick={onClickSave}>Save</AsyncOperationButton>
      </Grid>
    </Grid>
    <Grid item xs={12}>
      <StatusMessage {...status}/>
    </Grid>
  </SubPanel>
}

function MainIPSelect(props: { mainIp: string, addresses: string[], onMainIPSelected: (mainIp: string) => void }) {
  const disabled = props.addresses.length === 0
  const onChange = (e: ChangeEvent<HTMLSelectElement>) => props.onMainIPSelected(e.target.value)

  return <FormControl fullWidth>
    <InputLabel shrink={true} id="main-ip-label">Main IP address</InputLabel>
    <Select labelId="main-ip-label" value={props.mainIp} disabled={disabled} onChange={onChange}>
      {props.addresses.map(addr => <MenuItem key={addr} value={addr}>{addr}</MenuItem>)}
    </Select>
  </FormControl>
}

function StatePollingControls(props: { disabled: boolean, pollingEnabled: boolean, pollingIntervalSec: number, onPollChange: (enabled: boolean, interval: number) => void }) {
  const classes = useStyles()

  const onPollingEnabledChange = (e: ChangeEvent<HTMLInputElement>) => {
    props.onPollChange(e.target.checked, props.pollingIntervalSec)
  }
  const onPollIntervalChange = (e: ChangeEvent<HTMLSelectElement>) => {
    props.onPollChange(props.pollingEnabled, parseInt(e.target.value))
  }

  return <Grid item container xs={12} sm={9} md={10} lg={8} alignItems={'flex-end'}>
    <Grid item xs={12}>
      <Typography variant={'caption'} color={'textSecondary'}>State polling</Typography>
    </Grid>
    <Grid item container xs={12}>
      <Grid item container xs={5}>
        <Grid item xs={12}>
          <FormControlLabel
            control={
              <Checkbox disabled={props.disabled} checked={props.pollingEnabled} onChange={onPollingEnabledChange}
                        classes={{ root: classes.checkBoxRoot }}/>
            }
            label="Enabled"
          />
        </Grid>
      </Grid>
      <Grid item xs={3}>
        <FormControl fullWidth>
          <Select labelId="label" value={props.pollingIntervalSec}
                  disabled={props.disabled || !props.pollingEnabled} onChange={onPollIntervalChange}
                  endAdornment={
                    <InputAdornment position="start"
                                    className={classes.pollIntervalSelectAdornment}>min</InputAdornment>
                  }
          >
            {[1, 2, 3, 5, 10, 15, 20, 25, 30].map(v => <MenuItem key={v} value={v * 60}>{v}</MenuItem>)}
          </Select>
        </FormControl>
      </Grid>
    </Grid>
  </Grid>
}
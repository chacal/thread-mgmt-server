import { DeviceConfig, DeviceState } from './DeviceList'
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

const useStyles = makeStyles((theme) => ({
  configPanelRow: {
    marginBottom: theme.spacing(1)
  }
}))

export default function DeviceConfigPanel(props: { config: DeviceConfig, state: DeviceState, onSaveConfig: (c: DeviceConfig) => Promise<void> }) {
  const classes = useStyles()

  const [config, setConfig] = useState(props.config)
  const [status, setStatus] = useState(EmptyStatus)

  const onMainIPSelected = (mainIp: string) => setConfig({ ...config, mainIp })

  const onClickSave = () => {
    setStatus(EmptyStatus)
    return props.onSaveConfig(config)
      .catch(err => setStatus({ msg: err.toString(), isError: true, showProgress: false }))
  }

  const isSaveDisabled = () => isEqual(config, props.config)

  return <SubPanel heading={'Config'}>
    <Grid item container spacing={3} className={classes.configPanelRow}>
      <Grid item xs={12} sm={9} md={10} lg={8}>
        <MainIPSelect mainIp={config.mainIp} addresses={props.state.addresses} onMainIPSelected={onMainIPSelected}/>
      </Grid>
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

function MainIPSelect(props: { mainIp: string | undefined, addresses: string[] | undefined, onMainIPSelected: (mainIp: string) => void }) {
  const value = props.mainIp !== undefined ? props.mainIp : ''
  const disabled = props.addresses === undefined || props.addresses.length === 0
  const onChange = (e: ChangeEvent<HTMLSelectElement>) => props.onMainIPSelected(e.target.value)

  return <FormControl fullWidth>
    <InputLabel shrink={true} id="main-ip-label">Main IP address</InputLabel>
    <Select labelId="main-ip-label" value={value} disabled={disabled} onChange={onChange}>
      {
        props.addresses ? props.addresses.map(addr =>
          <MenuItem key={addr} value={addr}>{addr}</MenuItem>
        ) : null
      }
    </Select>
  </FormControl>
}


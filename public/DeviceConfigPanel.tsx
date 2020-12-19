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

  const onMainIPSelected = (e: ChangeEvent<HTMLSelectElement>) => {
    setConfig({ ...config, mainIp: e.target.value })
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
        <FormControl fullWidth>
          <InputLabel shrink={true} id="main-ip-label">Main IP address</InputLabel>
          <Select labelId="main-ip-label"
                  value={config.mainIp !== undefined ? config.mainIp : ''}
                  disabled={props.state.addresses === undefined || props.state.addresses.length === 0}
                  onChange={onMainIPSelected}
          >
            {
              props.state.addresses ? props.state.addresses.map(addr =>
                <MenuItem key={addr} value={addr}>{addr}</MenuItem>
              ) : null
            }
          </Select>
        </FormControl>
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
import React from 'react'
import { Device, DeviceConfig, DeviceDefaults, DeviceState } from './DeviceList'
import { makeStyles } from '@material-ui/core/styles'
import DeviceDefaultsPanel from './DeviceDefaultsPanel'
import Grid from '@material-ui/core/Grid'
import Paper from '@material-ui/core/Paper'
import DeviceStatePanel from './DeviceStatePanel'
import DeviceConfigPanel from './DeviceConfigPanel'
import DeviceTitleRow from './DeviceTitleRow'

const useStyles = makeStyles((theme) => ({
  root: {
    padding: theme.spacing(2, 4),
  }
}))

interface DeviceListItemProps {
  deviceId: string,
  device: Device,
  deviceChanged: (id: string, dev: Device) => void
  deviceRemoved: (id: string) => void
}

export default function DeviceListItem(props: DeviceListItemProps) {
  const classes = useStyles()

  const onSaveDefaults = (defaults: DeviceDefaults) => {
    const dev = { ...props.device, defaults }
    return postJSON('/v1/devices/' + props.deviceId + '/defaults', defaults)
      .then(() => props.deviceChanged(props.deviceId, dev))
  }

  const onSaveConfig = (config: DeviceConfig) => {
    const dev = { ...props.device, config }
    return postJSON('/v1/devices/' + props.deviceId + '/config', config)
      .then(() => props.deviceChanged(props.deviceId, dev))
  }

  const onStateRefresh = (state: DeviceState) => {
    const dev = { ...props.device, state }
    props.deviceChanged(props.deviceId, dev)
  }

  return <Grid item xs={12}>
    <Paper>
      <Grid container spacing={5} className={classes.root}>
        <DeviceTitleRow deviceId={props.deviceId} instance={props.device.defaults.instance}
                        deviceRemoved={props.deviceRemoved}/>
        <DeviceStatePanel state={props.device.state} deviceId={props.deviceId}
                          mainIp={props.device.config.mainIp} onStateRefresh={onStateRefresh}/>
        <DeviceDefaultsPanel defaults={props.device.defaults} deviceId={props.deviceId}
                             mainIp={props.device.config.mainIp} onSaveDefaults={onSaveDefaults}/>
        <DeviceConfigPanel config={props.device.config} addresses={props.device.state?.addresses} onSaveConfig={onSaveConfig}/>
      </Grid>
    </Paper>
  </Grid>
}

export function doFetch(input: RequestInfo, init?: RequestInit) {
  return fetch(input, init)
    .then(res => {
      if (res.status !== 200) {
        throw 'Status: ' + res.status
      }
      return res
    })
}

export function postJSON(url: string, data: any) {
  return doFetch(url, {
    method: 'POST',
    body: JSON.stringify(data),
    headers: { 'Content-Type': 'application/json' },
  })
}

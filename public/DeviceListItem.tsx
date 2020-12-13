import React from 'react'
import { Device, DeviceConfig, DeviceDefaults } from './DeviceList'
import { makeStyles } from '@material-ui/core/styles'
import DeviceDefaultsPanel from './DeviceDefaultsPanel'
import Grid from '@material-ui/core/Grid'
import Paper from '@material-ui/core/Paper'
import Typography from '@material-ui/core/Typography'
import DeviceStatePanel from './DeviceStatePanel'
import DeviceConfigPanel from './DeviceConfigPanel'

const useStyles = makeStyles((theme) => ({
  root: {
    padding: theme.spacing(2),
    paddingLeft: theme.spacing(4),
  },
  deviceId: {
    marginLeft: theme.spacing(2)
  }
}))

export default function DeviceListItem(props: { deviceId: string, device: Device, deviceSaved: (id: string, dev: Device) => void }) {
  const classes = useStyles()

  const onSaveDefaults = (defaults: DeviceDefaults) => {
    const dev = { ...props.device, defaults }
    return postJSON('/v1/devices/' + props.deviceId + '/defaults', defaults)
      .then(() => props.deviceSaved(props.deviceId, dev))
  }

  const onSaveConfig = (config: DeviceConfig) => {
    const dev = { ...props.device, config }
    return postJSON('/v1/devices/' + props.deviceId + '/config', config)
      .then(() => props.deviceSaved(props.deviceId, dev))
  }

  return <Grid item xs={12}>
    <Paper>
      <Grid container spacing={5} className={classes.root}>
        <TitleRow deviceId={props.deviceId} instance={props.device.defaults.instance}/>
        <DeviceStatePanel state={props.device.state}/>
        <DeviceDefaultsPanel defaults={props.device.defaults} onSaveDefaults={onSaveDefaults}/>
        <DeviceConfigPanel config={props.device.config} state={props.device.state} onSaveConfig={onSaveConfig}/>
      </Grid>
    </Paper>
  </Grid>
}

function TitleRow(props: { deviceId: string, instance?: string }) {
  const classes = useStyles()

  return <Grid item container alignItems={'flex-end'}>
    <Typography variant={'h5'} color={'primary'} display={'inline'}>
      {props.instance ? props.instance : 'N/A'}
    </Typography>
    <Typography variant={'subtitle1'} color={'textSecondary'} display={'inline'} className={classes.deviceId}>
      {props.deviceId}
    </Typography>
  </Grid>
}


function postJSON(url: string, data: any) {
  return fetch(url, {
    method: 'POST',
    body: JSON.stringify(data),
    headers: { 'Content-Type': 'application/json' },
  })
    .then(res => {
      if (res.status !== 200) {
        throw 'Status: ' + res.status
      }
    })
}

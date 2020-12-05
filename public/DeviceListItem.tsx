import React from 'react'
import { Device, DeviceAddress } from './DeviceList'
import { makeStyles } from '@material-ui/core/styles'
import DeviceSettingsPanel, { DeviceSettings } from './DeviceSettingsPanel'
import Grid from '@material-ui/core/Grid'
import Paper from '@material-ui/core/Paper'
import Typography from '@material-ui/core/Typography'
import IPAddressesPanel from './IPAddressesPanel'

const useStyles = makeStyles((theme) => ({
  root: {
    padding: theme.spacing(2),
    paddingLeft: theme.spacing(4),
  }
}))

export default function DeviceListItem(props: { deviceId: string, device: Device, deviceSaved: (id: string, dev: Device) => void }) {
  const classes = useStyles()

  const onSaveDevice = (dev: Device) => {
    return postJSON('/v1/devices/' + props.deviceId, dev)
      .then(() => props.deviceSaved(props.deviceId, dev))
  }

  const onSaveSettings = (s: DeviceSettings) => {
    const d = { ...props.device, ...s }
    return onSaveDevice(d)
  }

  const onSaveAddresses = (addr: DeviceAddress[]) => {
    const d = { ...props.device, addresses: addr }
    return onSaveDevice(d)
  }

  return <Grid item xs={12}>
    <Paper>
      <Grid container spacing={3} className={classes.root}>
        <TitleRow deviceId={props.deviceId} instance={props.device.instance}/>
        <IPAddressesPanel addresses={props.device.addresses} onSaveAddresses={onSaveAddresses}/>
        <DeviceSettingsPanel settings={settingsFor(props.device)} onSaveSettings={onSaveSettings}/>
      </Grid>
    </Paper>
  </Grid>
}

function settingsFor(d: Device): DeviceSettings {
  return {
    instance: d.instance,
    txPower: d.txPower,
    pollPeriod: d.pollPeriod
  }
}

function TitleRow(props: { deviceId: string, instance?: string }) {
  return <Grid item container alignItems={'flex-end'}>
    <Grid item xs={2} md={1}>
      <Typography variant={'h5'} color={'primary'}>
        {props.instance ? props.instance : 'N/A'}
      </Typography>
    </Grid>
    <Grid item xs={2} md={1}>
      <Typography variant={'subtitle1'} color={'textSecondary'}>
        {props.deviceId}
      </Typography>
    </Grid>
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

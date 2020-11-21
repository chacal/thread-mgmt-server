import React from 'react'
import { Device } from './DeviceList'
import { makeStyles } from '@material-ui/core/styles'
import DeviceSettingsPanel from './DeviceSettingsPanel'
import SubPanel from './SubPanel'
import Grid from '@material-ui/core/Grid'
import Paper from '@material-ui/core/Paper'
import Typography from '@material-ui/core/Typography'

const useStyles = makeStyles((theme) => ({
  root: {
    padding: theme.spacing(2),
    paddingLeft: theme.spacing(4),
  }
}))

export default function DeviceListItem(props: { deviceId: string, device: Device, deviceSaved: (id: string, dev: Device) => void }) {
  const classes = useStyles()

  const onSave = (dev: Device) => {
    return postJSON('/v1/devices/' + props.deviceId, dev)
      .then(() => props.deviceSaved(props.deviceId, dev))
  }

  return <Grid item xs={12}>
    <Paper>
      <Grid container xs={12} spacing={3} className={classes.root}>
        <TitleRow deviceId={props.deviceId} instance={props.device.instance}/>
        <IPAddressesPanel addresses={props.device.addresses}/>
        <DeviceSettingsPanel key={props.deviceId} dev={props.device} onSaveDevice={onSave}/>
      </Grid>
    </Paper>
  </Grid>
}

function TitleRow(props: { deviceId: string, instance: string }) {
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

function IPAddressesPanel(props: { addresses: string[] }) {
  return <SubPanel heading={'Addresses'}>
    {props.addresses ? props.addresses.map(addr =>
      <Grid item xs={12}>
        <Typography variant={'subtitle1'}>
          {addr}
        </Typography>
      </Grid>
    ) : null}
  </SubPanel>
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

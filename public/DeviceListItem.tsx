import React from 'react'
import { Grid, Paper, Typography } from '@material-ui/core'
import { Device } from './DeviceList'
import { makeStyles } from '@material-ui/core/styles'
import DeviceSettingsPanel from './DeviceSettingsPanel'
import SubPanel from './SubPanel'

const useStyles = makeStyles((theme) => ({
  root: {
    padding: theme.spacing(2),
    paddingLeft: theme.spacing(4),
  }
}))

export default function DeviceListItem(props: { deviceId: string, device: Device }) {
  const classes = useStyles()

  return <Grid item xs={12}>
    <Paper>
      <Grid container xs={12} spacing={3} className={classes.root}>
        <TitleRow deviceId={props.deviceId} instance={props.device.instance}/>
        <IPAddressesPanel addresses={props.device.addresses}/>
        <DeviceSettingsPanel device={props.device}/>
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
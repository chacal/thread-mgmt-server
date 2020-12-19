import React, { useState } from 'react'
import { Device, DeviceConfig, DeviceDefaults } from './DeviceList'
import { makeStyles } from '@material-ui/core/styles'
import DeviceDefaultsPanel from './DeviceDefaultsPanel'
import Grid from '@material-ui/core/Grid'
import Paper from '@material-ui/core/Paper'
import Typography from '@material-ui/core/Typography'
import DeviceStatePanel from './DeviceStatePanel'
import DeviceConfigPanel from './DeviceConfigPanel'
import DeleteIcon from '@material-ui/icons/DeleteOutline'
import IconButton from '@material-ui/core/IconButton'
import Snackbar from '@material-ui/core/Snackbar'
import Dialog from '@material-ui/core/Dialog'
import DialogActions from '@material-ui/core/DialogActions'
import DialogContent from '@material-ui/core/DialogContent'
import DialogContentText from '@material-ui/core/DialogContentText'
import DialogTitle from '@material-ui/core/DialogTitle'
import Button from '@material-ui/core/Button'
import Alert from '@material-ui/lab/Alert'

const useStyles = makeStyles((theme) => ({
  root: {
    padding: theme.spacing(2, 4),
  },
  deviceId: {
    marginLeft: theme.spacing(2)
  }
}))

interface DeviceListItemProps {
  deviceId: string,
  device: Device,
  deviceSaved: (id: string, dev: Device) => void
  deviceRemoved: (id: string) => void
}

export default function DeviceListItem(props: DeviceListItemProps) {
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
        <TitleRow deviceId={props.deviceId} instance={props.device.defaults.instance}
                  deviceRemoved={props.deviceRemoved}/>
        <DeviceStatePanel state={props.device.state}/>
        <DeviceDefaultsPanel defaults={props.device.defaults} deviceId={props.deviceId}
                             mainIp={props.device.config.mainIp} onSaveDefaults={onSaveDefaults}/>
        <DeviceConfigPanel config={props.device.config} state={props.device.state} onSaveConfig={onSaveConfig}/>
      </Grid>
    </Paper>
  </Grid>
}

function TitleRow(props: { deviceId: string, instance?: string, deviceRemoved: (deviceId: string) => void }) {
  const classes = useStyles()

  return <Grid item container xs={12}>
    <Grid item xs={8}>
      <Typography variant={'h5'} color={'primary'} display={'inline'}>
        {props.instance ? props.instance : 'N/A'}
      </Typography>
      <Typography variant={'subtitle1'} color={'textSecondary'} display={'inline'} className={classes.deviceId}>
        {props.deviceId}
      </Typography>
    </Grid>
    <Grid item container xs={4} justify={'flex-end'}>
      <DeleteButton deviceId={props.deviceId} deviceRemoved={props.deviceRemoved}/>
    </Grid>
  </Grid>
}

function DeleteButton(props: { deviceId: string, deviceRemoved: (deviceId: string) => void }) {
  const [confirmDialogOpen, setConfirmDialogOpen] = useState(false)
  const [inProgress, setInProgress] = useState(false)
  const [errorOpen, setErrorOpen] = useState(false)
  const [errorText, setErrorText] = useState('')

  const onClickRemove = () => setConfirmDialogOpen(true)
  const onCancel = () => setConfirmDialogOpen(false)
  const onConfirmRemove = () => {
    setConfirmDialogOpen(false)
    setInProgress(true)
    doFetch('/v1/devices/' + props.deviceId, { method: 'DELETE' })
      .then(() => props.deviceRemoved(props.deviceId))
      .catch(err => {
        setErrorText(err.toString())
        setErrorOpen(true)
      })
      .finally(() => setInProgress(false))
  }
  const onErrorClose = () => setErrorOpen(false)

  return <React.Fragment>
    <IconButton onClick={onClickRemove} disabled={inProgress}>
      <DeleteIcon/>
    </IconButton>
    <Dialog open={confirmDialogOpen} onClose={onCancel}>
      <DialogTitle>{'Confirm'}</DialogTitle>
      <DialogContent>
        <DialogContentText>Are you sure you want to remove device {props.deviceId}?</DialogContentText>
      </DialogContent>
      <DialogActions>
        <Button color="primary" onClick={onCancel}>Cancel</Button>
        <Button color="secondary" onClick={onConfirmRemove}>Remove</Button>
      </DialogActions>
    </Dialog>
    <Snackbar open={errorOpen} autoHideDuration={6000} onClose={onErrorClose} anchorOrigin={{
      vertical: 'bottom',
      horizontal: 'right',
    }}>
      <Alert onClose={onErrorClose} severity="error" variant={'filled'}>{errorText}</Alert>
    </Snackbar>
  </React.Fragment>
}

function doFetch(input: RequestInfo, init?: RequestInit) {
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

import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import React, { useState } from 'react'
import IconButton from '@material-ui/core/IconButton'
import DeleteIcon from '@material-ui/icons/DeleteOutline'
import Dialog from '@material-ui/core/Dialog'
import DialogTitle from '@material-ui/core/DialogTitle'
import DialogContent from '@material-ui/core/DialogContent'
import DialogContentText from '@material-ui/core/DialogContentText'
import DialogActions from '@material-ui/core/DialogActions'
import Button from '@material-ui/core/Button'
import Snackbar from '@material-ui/core/Snackbar'
import Alert from '@material-ui/lab/Alert'
import { makeStyles } from '@material-ui/core/styles'
import { doFetch } from './DeviceListItem'

const useStyles = makeStyles((theme) => ({
  deviceId: {
    marginLeft: theme.spacing(2)
  }
}))

export default function DeviceTitleRow(props: { deviceId: string, instance?: string, deviceRemoved: (deviceId: string) => void }) {
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
import Typography from '@material-ui/core/Typography'
import React from 'react'
import CircularProgress from '@material-ui/core/CircularProgress'
import { makeStyles } from '@material-ui/core/styles'

export const EmptyStatus: Status = { msg: '', isError: false, showProgress: false }

export interface Status {
  msg: string
  isError: boolean
  showProgress: boolean
}

const useStyles = makeStyles((theme) => ({
  main: {
    marginTop: '5px'
  },
  progress: {
    marginLeft: theme.spacing(1)
  },
}))

export default function StatusMessage(props: Status) {
  const classes = useStyles()
  const color = props.isError ? 'error' : 'textPrimary'

  return <Typography variant={'body2'} component={'div'} color={color} className={classes.main}>
    {props.msg}
    &nbsp;
    {props.showProgress ? <CircularProgress size={12} className={classes.progress}/> : null}
  </Typography>
}
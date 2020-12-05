import Typography from '@material-ui/core/Typography'
import React from 'react'
import { makeStyles } from '@material-ui/core/styles'

const useStyles = makeStyles((theme) => ({
  error: {
    marginLeft: theme.spacing(2)
  }
}))

export default function ErrorMessage(props: { msg: string }) {
  const classes = useStyles()

  return <Typography variant={'body1'} color={'error'} display={'inline'} className={classes.error}>
    {props.msg}
  </Typography>
}
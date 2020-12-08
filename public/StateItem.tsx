import React from 'react'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import { makeStyles } from '@material-ui/core/styles'

const useStyles = makeStyles((theme) => ({
  heading: {
    marginBottom: theme.spacing(0.5),
  },
}))


export default function StateItem(props: { heading: string, values?: string[] }) {
  const classes = useStyles()

  return <Grid item container>
    <Grid item className={classes.heading}>
      <Typography variant={'caption'} color={'textSecondary'} >{props.heading}</Typography>
    </Grid>
    <Grid container item>
      {props.values ? props.values.map(val =>
        <Grid item key={val} xs={12}>
          <Typography variant={'body2'}>{val}</Typography>
        </Grid>
      ) : null
      }
    </Grid>
  </Grid>
}

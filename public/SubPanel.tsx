import React from 'react'
import { Grid, Typography } from '@material-ui/core'

export default function SubPanel(props: { heading: string, children?: React.ReactNode }) {
  return <Grid item container direction={'column'} xs={12} md={6} spacing={1}>
    <Grid item>
      <Typography variant={'subtitle1'} color={'textSecondary'}>{props.heading}</Typography>
    </Grid>
    <Grid item container>
      {props.children}
    </Grid>
  </Grid>
}

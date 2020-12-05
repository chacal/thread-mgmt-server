import React from 'react'
import Typography from '@material-ui/core/Typography'
import Grid from '@material-ui/core/Grid'

export default function SubPanel(props: { heading: string, children?: React.ReactNode }) {
  return <Grid item container direction={'column'} xs={12} md={6} spacing={1}>
    <Grid item>
      <Typography variant={'h6'} color={'primary'}>{props.heading}</Typography>
    </Grid>
    <Grid item container>
      {props.children}
    </Grid>
  </Grid>
}

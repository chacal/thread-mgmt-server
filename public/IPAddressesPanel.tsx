import SubPanel from './SubPanel'
import React from 'react'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'

export default function IPAddressesPanel(props: { addresses?: string[] }) {
  return <SubPanel heading={'Addresses'}>
    {props.addresses ? props.addresses.map(addr =>
      <Grid item xs={12} key={addr}>
        <Typography variant={'body1'}>{addr}</Typography>
      </Grid>
    ) : null
    }
  </SubPanel>
}

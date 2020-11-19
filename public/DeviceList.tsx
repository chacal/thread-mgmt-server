import React, { useEffect, useState } from 'react'
import { Grid, Typography } from '@material-ui/core'
import { sortBy, toPairs } from 'lodash'
import DeviceListItem from './DeviceListItem'

interface Devices {
  [key: string]: Device
}

export interface Device {
  instance: string
  txPower: number
  pollPeriod: number
  addresses: string[]
}

export default function DeviceList() {
  const [devices, setDevices] = useState<Devices>({})

  useEffect(() => {
    loadDevices()
      .then(setDevices)
  }, [])

  return (
    <Grid container spacing={6}>
      {sortedDevices(devices).map(([deviceId, device]) => <DeviceListItem deviceId={deviceId} device={device}/>)}
    </Grid>
  )
}

function loadDevices() {
  return fetch(`/v1/devices`)
    .then(res => res.json())
}

function sortedDevices(devs: Devices) {
  return sortBy(toPairs(devs), ([id, d]) => d.instance, ([id, d]) => id)
}
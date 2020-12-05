import React, { useEffect, useState } from 'react'
import Grid from '@material-ui/core/Grid'
import sortBy from 'lodash/sortBy'
import toPairs from 'lodash/toPairs'
import DeviceListItem from './DeviceListItem'

interface Devices {
  [key: string]: Device
}

export interface DeviceAddress {
  ip: string
  main: boolean
}

export interface Device {
  instance?: string
  txPower?: number
  pollPeriod?: number
  addresses?: DeviceAddress[]
}

export default function DeviceList() {
  const [devices, setDevices] = useState<Devices>({})

  useEffect(() => {
    loadDevices()
      .then(setDevices)
  }, [])

  const deviceSaved = (deviceId: string, dev: Device) => {
    setDevices(prev => ({ ...prev, [deviceId]: dev }))
  }

  return (
    <Grid container spacing={7}>
      {
        sortedDevices(devices)
          .map(([deviceId, device]) =>
            <DeviceListItem key={deviceId} deviceId={deviceId} device={device} deviceSaved={deviceSaved}/>)
      }
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
import React, { useEffect, useState } from 'react'
import Grid from '@material-ui/core/Grid'
import sortBy from 'lodash/sortBy'
import toPairs from 'lodash/toPairs'
import DeviceListItem from './DeviceListItem'

interface Devices {
  [key: string]: Device
}

export interface DeviceDefaults {
  instance?: string
  txPower?: number
  pollPeriod?: number
}

interface ParentInfo {
  rloc16: string
  linkQualityIn: number
  linkQualityOut: number
  avgRssi: number
  latestRssi: number
}

export interface DeviceState {
  addresses?: string[]
  vcc?: number
  instance?: string
  parent?: ParentInfo
}

export interface DeviceConfig {
  mainIp?: string
}

export interface Device {
  defaults: DeviceDefaults
  state: DeviceState
  config: DeviceConfig
}

export default function DeviceList() {
  const [devices, setDevices] = useState<Devices>({})

  useEffect(() => {
    loadDevices()
      .then(setDevices)
  }, [])

  const deviceChanged = (deviceId: string, dev: Device) => {
    setDevices(prev => ({ ...prev, [deviceId]: dev }))
  }

  const deviceRemoved = (deviceId: string) => {
    setDevices(prev => {
      const { [deviceId]: id, ...rest } = prev
      return rest
    })
  }

  return (
    <Grid container spacing={7}>
      {
        sortedDevices(devices)
          .map(([deviceId, device]) =>
            <DeviceListItem key={deviceId} deviceId={deviceId} device={device}
                            deviceChanged={deviceChanged} deviceRemoved={deviceRemoved}/>)
      }
    </Grid>
  )
}

function loadDevices(): Promise<Devices> {
  return fetch(`/v1/devices`)
    .then(res => res.json())
}

function sortedDevices(devs: Devices) {
  return sortBy(toPairs(devs), ([id, d]) => d.defaults.instance, ([id, d]) => id)
}
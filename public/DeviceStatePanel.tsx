import { DeviceState } from './DeviceList'
import StateItem from './StateItem'
import React, { useState } from 'react'
import SubPanel from './SubPanel'
import Grid from '@material-ui/core/Grid'
import AsyncOperationButton from './AsyncOperationButton'
import StatusMessage, { EmptyStatus } from './StatusMessage'
import { postJSON } from './DeviceListItem'

interface DeviceStatePanelProps {
  state: DeviceState
  deviceId: string
  mainIp: string | undefined
  onStateRefresh: (s: DeviceState) => void
}

export default function DeviceStatePanel(props: DeviceStatePanelProps) {
  const [status, setStatus] = useState(EmptyStatus)

  const onClickRefresh = () => {
    setStatus({ msg: 'Refreshing state..', isError: false, showProgress: true })
    return postJSON(`/v1/devices/${props.deviceId}/refresh_state`, { address: props.mainIp })
      .then(res => res.json())
      .then(state => {
        setStatus(EmptyStatus)
        props.onStateRefresh(state)
      })
      .catch(err => setStatus({ msg: err.toString(), isError: true, showProgress: false }))
  }

  const isRefreshDisabled = () => props.mainIp === undefined

  return <SubPanel heading={'State'}>
    <Grid item container xs={12}>
      <Grid item xs={4}>
        <StateItem heading={'Instance'} values={props.state.instance}/>
      </Grid>
      <Grid item xs={4}>
        <StateItem heading={'Voltage'} values={voltageString(props.state)}/>
      </Grid>
      <Grid item xs={4}>
        <StateItem heading={'RLOC16'} values={props.state.parent?.rloc16}/>
      </Grid>
    </Grid>
    <Grid item container xs={12}>
      <Grid item xs={4}>
        <StateItem heading={'Link Quality In/Out'} values={linkQualityString(props.state)}/>
      </Grid>
      <Grid item xs={4}>
        <StateItem heading={'Latest RSSI'} values={rssiString(props.state.parent?.latestRssi)}/>
      </Grid>
      <Grid item xs={4}>
        <StateItem heading={'Avg RSSI'} values={rssiString(props.state.parent?.avgRssi)}/>
      </Grid>
    </Grid>
    <Grid item xs={12}>
      <StateItem heading={'Addresses'} values={props.state.addresses}/>
    </Grid>
    <Grid item container spacing={2} xs={12} alignItems={'center'}>
      <Grid item>
        <AsyncOperationButton disabled={isRefreshDisabled()} onClick={onClickRefresh}>Refresh</AsyncOperationButton>
      </Grid>
    </Grid>
    <Grid item xs={12}>
      <StatusMessage {...status}/>
    </Grid>
  </SubPanel>
}

function voltageString(s: DeviceState) {
  return s.vcc !== undefined ? (s.vcc / 1000).toFixed(3) + ' V' : ''
}

function linkQualityString(s: DeviceState) {
  return s.parent !== undefined ? `${s.parent.linkQualityIn}/${s.parent.linkQualityOut}` : ''
}

function rssiString(rssi: number | undefined) {
  return rssi !== undefined ? `${rssi}dBm` : ''
}
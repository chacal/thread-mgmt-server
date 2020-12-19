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
    <StateItem heading={'Addresses'} values={props.state.addresses}/>
    <StateItem heading={'Voltage'} values={[voltageString(props.state)]}/>
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
  return s.vcc ? s.vcc / 1000 + ' V' : ''
}
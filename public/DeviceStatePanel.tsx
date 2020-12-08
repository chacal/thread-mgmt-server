import { DeviceState } from './DeviceList'
import StateItem from './StateItem'
import React from 'react'
import SubPanel from './SubPanel'

export default function DeviceStatePanel(props: { state: DeviceState }) {
  return <SubPanel heading={'State'}>
    <StateItem heading={'Addresses'} values={props.state.addresses}/>
  </SubPanel>

}
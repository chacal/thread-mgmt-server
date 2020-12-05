import { DeviceAddress } from './DeviceList'
import SubPanel from './SubPanel'
import React, { useState } from 'react'
import Grid from '@material-ui/core/Grid'
import { makeStyles } from '@material-ui/core/styles'
import FormControl from '@material-ui/core/FormControl'
import FormControlLabel from '@material-ui/core/FormControlLabel'
import RadioGroup from '@material-ui/core/RadioGroup'
import AsyncOperationButton from './AsyncOperationButton'
import Radio from '@material-ui/core/Radio'
import ErrorMessage from './ErrorMessage'

const useStyles = makeStyles((theme) => ({
  saveButtonRow: {
    marginTop: theme.spacing(1),
  }
}))

export default function IPAddressesPanel(props: { addresses?: DeviceAddress[], onSaveAddresses: (addr: DeviceAddress[]) => Promise<void> }) {
  const classes = useStyles()

  const [mainAddr, setMainAddr] = useState(mainFromAddresses(props.addresses))
  const [saveError, setSaveError] = useState('')

  const onSave = () => {
    setSaveError('')
    const addresses = markMainAsTrue(props.addresses, mainAddr)
    return props.onSaveAddresses(addresses)
      .catch(err => setSaveError(err.toString()))
  }

  return <SubPanel heading={'Main Address'}>
    {(props.addresses === undefined || props.addresses.length === 0) ? null :
      <React.Fragment>
        <Grid item xs={12}>
          <FormControl component="fieldset">
            <RadioGroup value={mainAddr} onChange={event => setMainAddr(event.target.value)}>
              {props.addresses.map(addr =>
                <FormControlLabel key={addr.ip} value={addr.ip} control={<Radio/>} label={addr.ip}/>)
              }
            </RadioGroup>
          </FormControl>
        </Grid>
        <Grid item className={classes.saveButtonRow}>
          <AsyncOperationButton disabled={mainFromAddresses(props.addresses) === mainAddr} onClick={onSave}>
            Save
          </AsyncOperationButton>
          <ErrorMessage msg={saveError}/>
        </Grid>
      </React.Fragment>
    }
  </SubPanel>
}

function mainFromAddresses(addresses?: DeviceAddress[]): string {
  if (addresses === undefined) {
    return ''
  }

  const main = addresses.find(a => a.main)
  return main ? main.ip : ''
}

function markMainAsTrue(addresses: DeviceAddress[] | undefined, mainAddr: string) {
  const addrOrEmpty = addresses ? addresses : []
  const addr = addrOrEmpty.map(a => a.ip === mainAddr ? { ...a, main: true } : { ...a, main: false })
  return addr
}

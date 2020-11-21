import React, { useState } from 'react'
import { Button } from '@material-ui/core'

export default function AsyncOperationButton(props: { disabled: boolean, onClick: () => Promise<void>, children?: React.ReactNode }) {
  const [inProgress, setInProgress] = useState(false)

  return <Button variant={'outlined'} color={'primary'} disabled={props.disabled || inProgress}
                 onClick={() => {
                   setInProgress(true)
                   props.onClick()
                     .finally(() => setInProgress(false))
                 }}>
    {props.children}
  </Button>
}

import React from 'react'
import ReactDOM from 'react-dom'
import { AppBar, Container, Toolbar, Typography } from '@material-ui/core'
import DeviceList from './DeviceList'
import { makeStyles } from '@material-ui/core/styles'

const useStyles = makeStyles((theme) => ({
  mainContainer: {
    marginTop: theme.spacing(3),
  },
}))

function App() {
  const classes = useStyles()

  return (
    <React.Fragment>
      <AppBar>
        <Toolbar>
          <Typography variant="h6">
            Devices
          </Typography>
        </Toolbar>
      </AppBar>
      <Toolbar/>
      <Container maxWidth={'lg'} className={classes.mainContainer}>
        <DeviceList/>
      </Container>
    </React.Fragment>
  )
}

ReactDOM.render(<App/>, document.querySelector('#root'))
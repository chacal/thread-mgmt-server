import React from 'react'
import ReactDOM from 'react-dom'
import DeviceList from './DeviceList'
import { makeStyles } from '@material-ui/core/styles'
import AppBar from '@material-ui/core/AppBar'
import Typography from '@material-ui/core/Typography'
import Toolbar from '@material-ui/core/Toolbar'
import Container from '@material-ui/core/Container/Container'

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
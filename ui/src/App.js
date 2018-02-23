import React, { Component } from 'react';

import AppBar from 'material-ui/AppBar';
import Toolbar from 'material-ui/Toolbar';
import Typography from 'material-ui/Typography';
import Button from 'material-ui/Button';

import { Route, Redirect, withRouter } from "react-router-dom";

import Service from './Service';
import Login from './Login';

const rootStyle = {
  margin: 30,
  marginTop: 90,
}

const bgStyle = {
  position: 'fixed',
  top: 0,
  right: 0,
  bottom: 0,
  left: 0,
  background: '#eee',
  overflow: 'auto',
}

function checkAuth() {
  if (localStorage.getItem('logged')) {
    return true;
  }
  return false;
}

const PrivateRoute = ({ component: Component, ...props }) => {
  return <Route {...props} render={(props) => {
    return checkAuth() ? <Component {...props} /> : <Redirect to="/login" />
  }} />
}

class App extends Component {
  render() {

    return (
      <div style={bgStyle}>
        <AppBar position="fixed">
          <Toolbar>
            <Typography style={{ flex: 1 }} variant="title" color="inherit">
              Docker Manager
            </Typography>
            {/* {
              checkAuth() &&
              <Button color="inherit" onClick={() => {
                localStorage.removeItem('logged')
                this.props.history.push('/')
              }}>Logout</Button>
            } */}
          </Toolbar>
        </AppBar>

        <div style={rootStyle}>
          <Route exact path="/" component={Service} />
          {/* <Route path="/login" component={Login} /> */}
        </div>

      </div>
    );
  }
}

export default withRouter(App);

import React from 'react';
import LoginNfc from './LoginNfc';
import Login from './Login';

/*
 * LoginChooser
 * Will choose the right login page
 * depending of the presence of window.libnfc
 */
var LoginChooser = React.createClass({
  render() {
    return (
      <div className="container-fluid" >
        { !window.libnfc ? (
          <Login />
          ):(
          <LoginNfc />
          )}
      </div>
    );
  }
});

module.exports = LoginChooser;

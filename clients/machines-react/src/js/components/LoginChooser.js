import React from 'react';
import LoginNfc from './LoginNfc';
import Login from './Login';

var LoginChooser = React.createClass({
  render() {
    return (
      <div className="container-fluid" >
        { window.libnfc ? (
          <Login />
          ):(
          <LoginNfc />
          )}
      </div>
    );
  }
});

module.exports = LoginChooser;

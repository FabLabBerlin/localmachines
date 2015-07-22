import React from 'react';

var LoginNfc = React.createClass({
  render() {
    return (
      <form className="login-form" >
        <div className="nfc-login-info-icon">
          <i className="fa fa-credit-card" ></i>
        </div>
      </form>
    );
  }
});

module.exports = LoginNfc;

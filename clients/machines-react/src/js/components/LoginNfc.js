import React from 'react';
import {Navigation} from 'react-router';
import LoginStore from '../stores/LoginStore';
import LoginActions from '../actions/LoginActions';

/*
 * LoginNfc
 * Component to Login if there is window.libnfc
 */
var LoginNfc = React.createClass({

  mixins: [ Navigation ],

  errorNFCCallback(error) {
    window.libnfc.cardRead.disconnect(this.nfcLogin);
    window.libnfc.cardReaderError.disconnect(this.errorNFCCallback);
    toastr.error(error);
    setTimeout(this.getNFCUid, 2000);
  },

  nfcLogin(uid) {
    LoginActions.nfcLogin(uid);
  },

  getNFCUid() {
    window.libnfc.cardRead.connect(this.nfcLogin);
    window.libnfc.cardReaderError.connect(this.errorNFCCallback)
    window.libnfc.asyncScan();
  },

  onChangeLoginNFC() {
    if( LoginStore.getIsLogged() ) {
      this.replaceWith('/machine');
    } else {
      setTimeout(this.getNFCUid, 1000);
    }
  },

  componentWillUnmount() {
    window.libnfc.cardRead.disconnect(this.nfcLogin);
    window.libnfc.cardReaderError.disconnect(this.errorNFCCallback);
  },

  componentDidMount() {
    setTimeout(this.getNFCUid, 1000);
    LoginStore.onChangeLoginNFC = this.onChangeLogin;
  },

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

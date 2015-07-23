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

  /*
   * Callback called when nfc reader error occure
   */
  errorNFCCallback(error) {
    window.libnfc.cardRead.disconnect(this.nfcLogin);
    window.libnfc.cardReaderError.disconnect(this.errorNFCCallback);
    toastr.error(error);
    setTimeout(this.connectJsToQt, 2000);
  },

  /*
   * Called if uid is received by the nfcReader
   */
  nfcLogin(uid) {
    LoginActions.nfcLogin(uid);
  },

  /*
   * Connect nfc event to javascript function
   * start the nfc polling
   */
  connectJsToQt() {
    window.libnfc.cardRead.connect(this.nfcLogin);
    window.libnfc.cardReaderError.connect(this.errorNFCCallback)
    window.libnfc.asyncScan();
  },

  /*
   * When LoginStore is done with his work
   */
  onChangeLoginNFC() {
    if( LoginStore.getIsLogged() ) {
      this.replaceWith('/machine');
    } else {
      setTimeout(this.connectJsToQt, 1000);
    }
  },

  /*
   * Destructor
   * disconnect the listener
   */
  componentWillUnmount() {
    window.libnfc.cardRead.disconnect(this.nfcLogin);
    window.libnfc.cardReaderError.disconnect(this.errorNFCCallback);
  },

  /*
   * Connect the onChangeLogin together
   * Start the connection between Qt and JS
   */
  componentDidMount() {
    setTimeout(this.connectJsToQt, 1000);
    LoginStore.onChangeLoginNFC = this.onChangeLogin;
  },

  /*
   * Render the big creditCard icon
   */
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

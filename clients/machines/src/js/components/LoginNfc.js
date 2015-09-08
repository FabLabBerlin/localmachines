var $ = require('jquery');
var getters = require('../getters');
var LoginStore = require('../stores/LoginStore');
var LoginActions = require('../actions/LoginActions');
var {Navigation} = require('react-router');
var React = require('react');
var reactor = require('../reactor');
var toastr = require('toastr');
var actionTypes = require('../actionTypes');

/*
 * LoginNfc
 * Component to Login if there is window.libnfc
 */
var LoginNfc = React.createClass({

  /*
   * To use transitionTo/replaceWith/redirect and some function related to the router
   */
  mixins: [ Navigation ],

  /*
   * Callback called when nfc reader error occure
   */
  errorNFCCallback(error) {
    console.log('errorNFCCallback');

    try {
      window.libnfc.cardRead.disconnect(this.nfcLogin);
      window.libnfc.cardReaderError.disconnect(this.errorNFCCallback);
    } catch (e) {
      console.log(e.message);
    }

    setTimeout(this.connectJsToQt, 1000);
  },

  /*
   * Called if uid is received by the nfcReader
   */
  nfcLogin(uid) {
    try {
      window.libnfc.cardRead.disconnect(this.nfcLogin);
      window.libnfc.cardReaderError.disconnect(this.errorNFCCallback);
    } catch (e) {
      console.log(e.message);
    }

    LoginActions.nfcLogin(uid, this.context.router);
  },

  /*
   * Connect nfc event to javascript function
   * start the nfc polling
   */
  connectJsToQt() {
    //toastr.info('connectJsToQt');
    try {
      window.libnfc.cardRead.connect(this.nfcLogin);
      window.libnfc.cardReaderError.connect(this.errorNFCCallback);
      window.libnfc.asyncScan();
    } catch (e) {
      console.log(e.message);
    }
  },

  /*
   * When LoginStore is done with his work
   */
  onChangeLoginNFC() {
    const isLogged = reactor.evaluateToJS(getters.getIsLogged);
    if (isLogged) {
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
    try {
      window.libnfc.cardRead.disconnect(this.nfcLogin);
      window.libnfc.cardReaderError.disconnect(this.errorNFCCallback);
    } catch (e) {
      // TODO: redirect the message to futer remote API
      console.log(e.message);
    }
  },

  /*
   * Connect the onChangeLogin together
   * Start the connection between Qt and JS
   */
  componentDidMount() {

    // For debugging through the console
    window.nfcLogin = this.nfcLogin;

    reactor.reset();

    setTimeout(this.connectJsToQt, 1000);
  },

  /*
   * Render NFC terminal welcome screen with info
   * about ways to log in to the system.
   */
  render() {
    return (


        <div className="row">

          <div className="col-xs-1"></div>

          <div className="col-xs-4">
            <div className="nfc-login-icon">
              <i className="fa fa-credit-card"></i>
            </div>
            <div className="nfc-login-info-text">
              <p>Use your NFC card to log in</p>
            </div>
          </div>

          <div className="col-xs-1">
            <div className="nfc-login-info-text">
              <p>or</p>
            </div>
          </div>

          <div className="col-xs-5">
            <div className="nfc-login-icon">
              <i className="fa fa-laptop"></i>
            </div>
            <div className="nfc-login-info-text">
              <p>go to <b>easylab.io</b> on your device.</p>
            </div>
          </div>

          <div className="col-xs-1"></div>

        </div>


    );
  }
});

export default LoginNfc;

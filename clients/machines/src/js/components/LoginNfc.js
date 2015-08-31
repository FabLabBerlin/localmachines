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

    LoginActions.nfcLogin(uid);
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

    setTimeout(this.connectJsToQt, 1000);

    reactor.reset(); // Otherwise we get observer clones

    reactor.observe(getters.getIsLogged, isLogged => {
      console.log('isLogged observer');
      this.onChangeLoginNFC();
    }.bind(this));

    reactor.observe(getters.getLoginSuccess, loginSuccess => {
      console.log('loginSuccess observer');
      //const loginFailure = reactor.evaluateToJS(getters.getLoginFailure);
      console.log('loginSuccess: ' + loginSuccess);
      if (!loginSuccess) {
        setTimeout(function() {
          this.connectJsToQt();
          reactor.dispatch(actionTypes.LOGIN_FAILURE_HANDLED);
        }.bind(this), 1000);
      }
    }.bind(this));
  },

  /*
   * Render the big creditCard icon
   */
  render() {
    return (
      <div>
        <h3 id="login-own-device-text">You can use EasyLab on your own Smartphone or Laptop!</h3>
        <div className="row">
          <div className="col-xs-5">
            <form className="login-form">
              <div className="nfc-login-info-icon">
                <i className="fa fa-credit-card"></i>
              </div>
              <div className="nfc-login-info-text">
                <p>Use your NFC card to log in</p>
              </div>
            </form>
          </div>
          <div id="login-nfc-or-own" className="col-xs-2">
            or
          </div>
          <div className="col-xs-5 login-own-device">
            <div id="login-laptop-overlay">
              <div>Username</div>
              <div>Password</div>
            </div>
            <img className="login-laptop" src="/machines/img/Laptop1.svg"/>
            <img className="login-phone" src="/machines/img/Phone1.svg"/>
            <div>
              Go to https://easylab.io on your Laptop or Phone
            </div>
          </div>
        </div>
      </div>
    );
  }
});

export default LoginNfc;

var getters = require('../getters');
var reactor = require('../reactor');
var toastr = require('../toastr');


var NfcLogoutMixin = {

  nfcOnDidMount() {
    // For debugging through the console
    window.nfcLogout = this.handleLogout;

    if(window.libnfc) {
      setTimeout(this.connectJsToQt, 1500);
    }

    reactor.observe(getters.getIsLogged, isLogged => {
      this.onChangeLogout();
    }.bind(this));
  },

  nfcOnWillUnmount() {
    if(window.libnfc){
      window.libnfc.cardRead.disconnect(this.handleLogout);
      window.libnfc.cardReaderError.disconnect(this.errorNFCCallback);
    }
  },

  /*
   * To logout and redirect to login page
   */
  onChangeLogout() {
    const isLogged = reactor.evaluateToJS(getters.getIsLogged);
    if (!isLogged) {
      this.replaceWith('login');
    }
  },

  connectJsToQt() {
    if (window.location.hash === '#/machine') {
      toastr.info('You can log out with your nfc card');
    }
    window.libnfc.cardRead.connect(this.handleLogout);
    window.libnfc.cardReaderError.connect(this.errorNFCCallback);
    window.libnfc.asyncScan();
  },

  /*
   * Callback called when nfc reader error occure
   */
  errorNFCCallback(error) {
    window.libnfc.cardRead.disconnect(this.nfcLogin);
    window.libnfc.cardReaderError.disconnect(this.errorNFCCallback);
    toastr.error(error);
    setTimeout(this.connectJsToQt, 2000);
  }

};

export default NfcLogoutMixin;

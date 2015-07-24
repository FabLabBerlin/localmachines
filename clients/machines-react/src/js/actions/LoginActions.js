import LoginStore from '../stores/LoginStore';
import toastr from 'toastr';
toastr.options.positionClass = 'toast-bottom-left';

/*
 * Action made by the login page
 */
var LoginActions = {

  /*
   * Submit login form to log in
   */
  submitLoginForm(content) {
    LoginStore.apiPostLogin(content);
  },

  /*
   * Try to connect with nfc card
   * @uid: unique id from the card
   */
  nfcLogin(uid) {
    LoginStore.apiPostLoginNFC(uid);
  },

  /*
   * Logout
   */
  logout() {
    LoginStore.apiGetLogout();
  }
};

module.exports = LoginActions;

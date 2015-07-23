import LoginStore from '../stores/LoginStore';

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
    console.log(uid);
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

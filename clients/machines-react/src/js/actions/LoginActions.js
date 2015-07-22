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
   * Logout
   */
  logout() {
    LoginStore.apiGetLogout();
  }
};

module.exports = LoginActions;

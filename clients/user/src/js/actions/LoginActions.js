import UserStore from '../stores/UserStore';

/*
 * All the actions called by the LoginPage
 */
var LoginActions = {

  /*
   * Try to login with the information in the form
   * @content: The login page form content
   * call the UserPage to interact with the back-end
   */
    submitLoginForm(content) {
        UserStore.submitLoginFormToServer(content);
    }

};

module.exports = LoginActions;

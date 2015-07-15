import UserStore from '../stores/UserStore'

/*
 * All the actions called by the UserPage
 */
var UserActions = {

  /*
   * Try to update the user information
   * @userState: data from userForm
   * call the UserStore to interact with the back-end
   */
  submitState(userState){
    UserStore.submitUpdatedStateToServer(userState);
  },

  /*
   * Ask the store to update the password
   * @password: new password the user want to have
   */
  updatePassword(password) {
    UserStore.updatePassword(password);
  },

  /*
   * To logout
   */
  logout() {
    UserStore.logoutFromServer();
  }

};

module.exports = UserActions;

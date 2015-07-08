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
   * To logout
   */
  logout() {
    UserStore.logoutFromServer();
  }

};

module.exports = UserActions;

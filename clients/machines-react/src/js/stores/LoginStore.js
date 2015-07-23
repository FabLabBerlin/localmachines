
/*
 * import toastr and set position
 */
import toastr from 'toastr';
toastr.options.positionClass = 'toast-bottom-left';

/*
 * Login Store:
 * The goal of this file is to handle all login related actions
 * state
 * apitGetLogout
 * apitPostLogin
 * getter
 * cleanState
 * onChange
 */
var LoginStore = {
  state: {
    firstTry: true,
    isLogged: false,
    uid: {}
  },

  /*
   * To logout
   */
  apiGetLogout() {
    $.ajax({
      url: '/api/users/logout',
      type: 'GET',
      cache: false,
      success: function(data) {
        this.cleanState();
      }.bind(this),
      error: function(xhr, status, err) {
        console.error('/users/logout', status, err);
      }.bind(this),
    });
  },

  /*
   * To login
   */
  apiPostLogin(loginInfo) {
    $.ajax({
      url: '/api/users/login',
      dataType: 'json',
      type: 'POST',
      data: loginInfo,
      success: function(data) {
        this.successLogin(data);
      }.bind(this),
      error: function(xhr, status, err) {
        this.errorLogin();
        console.error('/users/login', status, err);
      }.bind(this),
    });
  },

  apitPostLoginNFC(uid) {
    $.ajax({
      url: '/api/users/loginuid',
      method: 'POST',
      data: {
        uid: uid
      },
      success: function(data) {
        this.successLogin(data);
      }.bind(this),
      error: function(xhr, status, err) {
        this.errorLogin();
        console.error('/users/loginuid', status, err);
      }.bind(this),
    });
  },

  successLogin(data) {
    if( data.UserId ) {
      LoginStore.state.uid = data.UserId;
      LoginStore.putLoginState();
    } else {
      toastr.error('Failed to log in');
      LoginStore.putLoginState(false);
    }
  },

  errorLogin() {
    if(LoginStore.state.firstTry === true) {
      LoginStore.state.firstTry = false;
    } else {
      toastr.error('Wrong password');
    }
  },

  /*
   * Return the uid of the user
   */
  getUid() {
    return this.state.uid;
  },

  getIsLogged() {
    return this.state.isLogged;
  },

  cleanState() {
    this.state.isLogged = false;
    this.state.userInfo = {};
    console.log(toastr.success);
    toastr.success('Bye');
    this.onChangeLogout();
  },

  /*
   * Change state before login
   */
  putLoginState(log = true) {
    if( log === true ) {
      this.state.isLogged = true;
      this.state.firstTry = true;
    }
    this.onChangeLogin();
  },

  /*
   * Event triggered when login
   * See Login page
   */
  onChangeLogin() {},

  onChangeLogout() {}

}

module.exports = LoginStore;

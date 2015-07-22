
/*
 * import toastr and set position
 */
import toastr from 'toastr';
toastr.options.positionclass = 'toast-bottom-left';

var LoginStore = {
  state: {
    firstTry: true,
    isLogged: false,
    userInfo: {}
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
        this.state.userInfo.UserId = data.UserId;
        this.putLoginState();
      }.bind(this),
      error: function(xhr, status, err) {
        if(this.state.firstTry === true) {
          this.state.firstTry = false;
        } else {
          toastr.error('Wrong password');
        }
        console.error('/users/login', status, err);
      }.bind(this),
    });
  },

  /*
   * Return the uid of the user
   */
  getUid() {
    return this.state.userInfo.UserId;
  },

  getIsLogged() {
    return this.state.isLogged;
  },

  cleanState() {
    this.state.isLogged = false;
    this.state.userInfo = {};
    toastr.success('Bye');
    this.onChangeLogout();
  },

  /*
   * Change state before login
   */
  putLoginState() {
    this.state.isLogged = true;
    this.state.firstTry = true;
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

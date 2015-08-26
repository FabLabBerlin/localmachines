var $ = require('jquery');
var actionTypes = require('../actionTypes');
var reactor = require('../reactor');
var toastr = require('../toastr');


/*
 * Action made by the login page
 */
export default {

  /*
   * Submit login form to log in
   */
  submitLoginForm(content) {
    $.ajax({
      url: '/api/users/login',
      dataType: 'json',
      type: 'POST',
      data: content,
      success: function(data) {
        reactor.dispatch(actionTypes.SUCCESS_LOGIN, { data });
      }.bind(this),
      error: function(xhr, status, err) {
        if (content.username !== '' && content.password !== '') {
          toastr.error('Failed to log in');
        }
        reactor.dispatch(actionTypes.ERROR_LOGIN);
        console.error('/users/login', status, err);
      }.bind(this)
    });
  },

  /*
   * Try to connect with nfc card
   * @uid: unique id from the card
   */
  nfcLogin(uid) {
    $.ajax({
      url: '/api/users/loginuid',
      method: 'POST',
      data: {
        uid: uid
      },
      success: function(data) {
        reactor.dispatch(actionTypes.SUCCESS_LOGIN, { data });
      }.bind(this),
      error: function(xhr, status, err) {
        reactor.dispatch(actionTypes.ERROR_LOGIN);
        //console.error('/users/loginuid', status, err);
      }.bind(this)
    });
  },

  keepAlive() {
    reactor.dispatch(actionTypes.KEEP_ALIVE);
  },

  /*
   * Logout
   */
  logout() {
    $.ajax({
      url: '/api/users/logout',
      type: 'GET',
      cache: false,
      success: function(data) {
        reactor.dispatch(actionTypes.SUCCESS_LOGOUT);
      }.bind(this),
      error: function(xhr, status, err) {
        console.error('/users/logout', status, err);
      }.bind(this)
    });
  }

};

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
  submitLoginForm(content, router) {
    $.ajax({
      url: '/api/users/login',
      dataType: 'json',
      type: 'POST',
      data: content,
      success(data) {
        reactor.dispatch(actionTypes.SUCCESS_LOGIN, { data });
        router.transitionTo('/machine');
      },
      error(xhr, status, err) {
        if (content.username !== '' && content.password !== '') {
          toastr.error('Failed to log in');
        }
        reactor.dispatch(actionTypes.ERROR_LOGIN);
        console.error('/users/login', status, err);
      }
    });
  },

  /*
   * Try to connect with nfc card
   * @uid: unique id from the card
   */
  nfcLogin(uid, router) {
    $.ajax({
      url: '/api/users/loginuid',
      method: 'POST',
      data: {
        uid: uid
      },
      success(data) {
        reactor.dispatch(actionTypes.SUCCESS_LOGIN, { data });
        router.transitionTo('/machine');
      },
      error(xhr, status, err) {
        toastr.error('Problem with NFC login.  Try again later or talk to us if the problem persists.');
        setTimeout(function() {
          document.location.reload(true);
        }, 2000);
        reactor.dispatch(actionTypes.ERROR_LOGIN);
        //console.error('/users/loginuid', status, err);
      }
    });
  },

  keepAlive() {
    reactor.dispatch(actionTypes.KEEP_ALIVE);
  },

  /*
   * Logout
   */
  logout(router) {
    $.ajax({
      url: '/api/users/logout',
      type: 'GET',
      cache: false,
      success(data) {
        reactor.dispatch(actionTypes.SUCCESS_LOGOUT);
        if (router) {
          router.transitionTo('/login');
        } else {
          toastr.info('router not defined');
        }
      },
      error(xhr, status, err) {
        console.error('/users/logout', status, err);
      }
    });
  }

};

var $ = require('jquery');
var actionTypes = require('../actionTypes');
var LocationActions = require('./LocationActions');
var Machines = require('../modules/Machines');
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
        LocationActions.setLocationId(data.LocationId);
        reactor.dispatch(actionTypes.SUCCESS_LOGIN, { data });
        switch (data.Status) {
        case 'ok':
          router.transitionTo('/machine');
          break;
        case 'unregistered':
          router.transitionTo('/register_existing');
          break;
        default:
          console.log('unknowon status');
          toastr.error('Error.  Please try again later.');
        }
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
   * Maybe the user is already logged in
   */
  tryPassLoginForm(router, cb) {
    $.ajax({
      url: '/api/users/current',
      dataType: 'json',
      type: 'GET',
      params: {
        ac: new Date().getTime()
      },
      success(user) {
        if (user) {
          var data = {
            UserId: user.Id
          };
          reactor.dispatch(actionTypes.SUCCESS_LOGIN, { data });
        } else {
          if (cb) {
            cb();
          }
        }
      },
      error(xhr, status, err) {
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
        uid: uid,
        location: 1 // Hardcoded for now
      },
      success(data) {
        reactor.dispatch(actionTypes.SUCCESS_LOGIN, { data });
        LocationActions.setLocationId(1);
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
    //reactor.dispatch(actionTypes.KEEP_ALIVE);
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
        window.location.href = '/machines/#/login';
      },
      error(xhr, status, err) {
        console.error('/users/logout', status, err);
      }
    });
  }

};

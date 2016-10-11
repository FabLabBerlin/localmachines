var $ = require('jquery');
var actionTypes = require('../actionTypes');
var LocationActions = require('./LocationActions');
var Machines = require('../modules/Machines');
var reactor = require('../reactor');
var toastr = require('../toastr');

import {hashHistory} from 'react-router';


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
          hashHistory.push('/machines');
          break;
        case 'unregistered':
          hashHistory.push('/register_existing');
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
  tryAutoLogin(router, {loggedIn, loggedOut}) {
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
          reactor.dispatch(actionTypes.SUCCESS_AUTO_LOGIN, { data });
          if (loggedIn) {
            loggedIn();
          }
        } else {
          reactor.dispatch(actionTypes.FAIL_AUTO_LOGIN);
          if (loggedOut) {
            loggedOut();
          }
        }
      },
      error(xhr, status, err) {
        reactor.dispatch(actionTypes.FAIL_AUTO_LOGIN);
        console.error('/users/login', status, err);
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

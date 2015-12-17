var $ = require('jquery');
var actionTypes = require('../actionTypes');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;
var toastr = require('../toastr');


const initialState = toImmutable({
  firstTry: true,
  isLogged: false,
  loginSuccess: true,
  uid: {},
  lastActivity: new Date()
});

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
var LoginStore = new Nuclear.Store({
  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.SUCCESS_LOGIN, successLogin);
    this.on(actionTypes.ERROR_LOGIN, errorLogin);
    this.on(actionTypes.SUCCESS_LOGOUT, successLogout);
    this.on(actionTypes.LOGIN_FAILURE_HANDLED, onLoginFailureHandled);
    this.on(actionTypes.KEEP_ALIVE, keepAlive);
  }
});

/*
 * Success to the login functions
 * if data is corrupted, say it to putLoginState
 */
function successLogin(state, { data }) {
  if (data.UserId) {
    return putLoginState(state.set('uid', data.UserId));
  } else {
    toastr.error('Failed to log in');
    return putLoginState(state, false);
  }
}

/*
 * Error callback of the login functions
 */
function errorLogin(state) {
  if (state.get('firstTry')) {
    return state.set('firstTry', false)
                .set('loginSuccess', false);
  } else {
    return state.set('loginSuccess', false);
  }
}

/*
 * Return the uid of the user
 */
function getUid() {
  return this.state.uid;
}

function getIsLogged() {
  return this.state.isLogged;
}

/*
 * Clean the store before logout
 */
function successLogout(state) {
  toastr.success('Bye');
  onChangeLogout();
  return state.set('isLogged', false)
              .set('user', {});
}

function onLoginFailureHandled(state) {
  //console.log('onLoginFailureHandled');
  return state.set('loginSuccess', true);
}

/*
 * Change state before login
 * If fail to log, don't change the store
 */
function putLoginState(state, log = true) {
  onChangeLogin();
  state = keepAlive(state);
  if (log) {
    return state.set('isLogged', true)
                .set('firstTry', true);
  } else {
    return state;
  }
}

function keepAlive(state) {
  return state.set('lastActivity', new Date());
}

/*
 * Event triggered when login
 * See Login page
 */
function onChangeLogin() {}

function onChangeLogout() {}

export default LoginStore;

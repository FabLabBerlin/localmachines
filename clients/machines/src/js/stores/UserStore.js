var $ = require('jquery');
var actionTypes = require('../actionTypes');
var Nuclear = require('nuclear-js');
var reactor = require('../reactor');
var toastr = require('../toastr');
var toImmutable = Nuclear.toImmutable;


const initialState = toImmutable({
  userId: 0,
  isLogged: false,
  firstTry: true,
  billInfo: undefined,
  membershipInfo: [],
  userInfo: {}
});

/*
 * @UserStore:
 * All the information about the user are stored here
 * All the interaction between the front-end and the back-end are done here
 * @state
 * @ajax calls:
 *  - PUT
 *  - POST
 *  - GET
 * @formatfunction
 * @getter
 * TODO: It would be possible to have multiple membership
 *       or to keep trace of your spending in membership
 *       Need to Change the reponse of /users/uid/membership by an array
 *
 */
var UserStore = new Nuclear.Store({
  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.SET_MEMBERSHIP_INFO, setMembershipInfo);
    this.on(actionTypes.SET_BILL_INFO, setBillInfo);
    this.on(actionTypes.SET_USER_INFO, setUserInfo);
    this.on(actionTypes.SET_USER_INFO_PROPERTY, setUserInfoProperty);
  }

});

function setMembershipInfo(state, { data }) {
  return state.set('membershipInfo', data);
}

function setBillInfo(state, { data }) {
  return state.set('billInfo', data);
}

function setUserInfo(state, { userInfo }) {
  return state.set('userInfo', toImmutable(userInfo));
}

function setUserInfoProperty(state, { key, value }) {
  return state.set('userInfo', state.get('userInfo').set(key, value));
}

module.exports = UserStore;

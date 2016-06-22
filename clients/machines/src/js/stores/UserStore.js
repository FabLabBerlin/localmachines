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
  bill: undefined,
  memberships: undefined,
  user: {}
});

var UserStore = new Nuclear.Store({
  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.SET_MEMBERSHIPS, setMemberships);
    this.on(actionTypes.SET_BILL, setBill);
    this.on(actionTypes.SET_USER, setUser);
    this.on(actionTypes.SET_USER_PROPERTY, setUserProperty);
  }

});

function setMemberships(state, { data }) {
  return state.set('memberships', toImmutable(data));
}

function setBill(state, { data }) {
  return state.set('bill', toImmutable(data));
}

function setUser(state, { user }) {
  return state.set('user', toImmutable(user));
}

function setUserProperty(state, { key, value }) {
  return state.set('user', state.get('user').set(key, value));
}

export default UserStore;

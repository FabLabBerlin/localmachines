var actionTypes = require('../actionTypes');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;


const initialState = toImmutable({
  userMemberships: {},
  users: undefined
});


var UsersStore = new Nuclear.Store({
  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.SET_USER_MEMBERSHIPS, setUserMemberships);
    this.on(actionTypes.SET_USERS, setUsers);
  }
});

function setUserMemberships(state, {userId, userMemberships}) {
  return state.setIn(['userMemberships', userId], userMemberships);
}

function setUsers(state, users) {
  return state.set('users', users);
}

export default UsersStore;

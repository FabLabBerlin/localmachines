var actionTypes = require('../actionTypes');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;


const initialState = toImmutable({});


var UsersStore = new Nuclear.Store({
  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.SET_USERS, setUsers);
  }
});

function setUsers(state, users) {
  return state.set('users', users);
}

export default UsersStore;

var actionTypes = require('../actionTypes');
var Nuclear = require('nuclear-js');
var reactor = require('../../../reactor');
var toImmutable = Nuclear.toImmutable;


const initialState = toImmutable({});

var ForgotPasswordStore = new Nuclear.Store({
  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.SET_KEY, setKey);
    this.on(actionTypes.SET_PHONE, setPhone);
  }
});

function setKey(state, key) {
  return state.set('key', key);
}

function setPhone(state, phone) {
  return state.set('phone', phone);
}

export default ForgotPasswordStore;

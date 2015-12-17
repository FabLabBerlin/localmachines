var actionTypes = require('../actionTypes');
var Nuclear = require('nuclear-js');
var reactor = require('../reactor');
var toastr = require('../toastr');
var toImmutable = Nuclear.toImmutable;


const initialState = toImmutable({
  loading: false
});

var GlobalStore = new Nuclear.Store({
  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.SET_LOADING, setLoading);
    this.on(actionTypes.UNSET_LOADING, unsetLoading);
  }
});

function setLoading(state) {
  return state.set('loading', true);
}

function unsetLoading(state) {
  return state.set('loading', false);
}

export default GlobalStore;

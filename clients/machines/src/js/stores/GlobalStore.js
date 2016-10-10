var actionTypes = require('../actionTypes');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;


const initialState = toImmutable({
  loading: false,
  width: window.innerWidth,
  height: window.innerHeight
});

var GlobalStore = new Nuclear.Store({
  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.SET_WINDOW_SIZE, setWindowSize);
    this.on(actionTypes.SET_LOADING, setLoading);
    this.on(actionTypes.UNSET_LOADING, unsetLoading);
  }
});

function setWindowSize(state, {width, height}) {
  return state.set('width', width)
              .set('height', height);
}

function setLoading(state) {
  return state.set('loading', true);
}

function unsetLoading(state) {
  return state.set('loading', false);
}

export default GlobalStore;

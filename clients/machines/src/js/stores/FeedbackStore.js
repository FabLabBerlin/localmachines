var $ = require('jquery');
var actionTypes = require('../actionTypes');
var Nuclear = require('nuclear-js');
var reactor = require('../reactor');
var toastr = require('../toastr');
var toImmutable = Nuclear.toImmutable;


const initialState = toImmutable({
  ['subject-dropdown']: 'Billing',
  ['subject-other-text']: '',
  message: ''
});

var FeedbackStore = new Nuclear.Store({
  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.SET_FEEDBACK_PROPERTY, setProperty);
    this.on(actionTypes.RESET_FEEDBACK_FORM, reset);
  }
});

function setProperty(state, { key, value }) {
  return state.set(key, value);
}

function reset(state) {
  return initialState;
}

export default FeedbackStore;

var $ = require('jquery');
var actionTypes = require('../actionTypes');
var Nuclear = require('nuclear-js');
var reactor = require('../reactor');
var toastr = require('../toastr');
var toImmutable = Nuclear.toImmutable;


const initialState = toImmutable({
  tutorings: undefined
});

var TutoringsStore = new Nuclear.Store({
  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.SET_TUTORINGS, setTutorings);
  }
});

function setTutorings(state, tutorings) {
  return state.set('tutorings', toImmutable(tutorings));
}

module.exports = TutoringsStore;

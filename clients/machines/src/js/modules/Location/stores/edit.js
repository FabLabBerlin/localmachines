import _ from 'lodash';
import actionTypes from '../actionTypes';
import Nuclear from 'nuclear-js';
var toImmutable = Nuclear.toImmutable;

const initialState = toImmutable({
  Id: undefined
});

var LocationEditStore = new Nuclear.Store({
  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.SET_EDIT_LOCATION, setEditLocation);
  }
});

function setEditLocation(state, location) {
  var s = state;

  _.each(location, (v, k) => {
    s = s.set(k, v);
  });

  return s;
}

export default LocationEditStore;

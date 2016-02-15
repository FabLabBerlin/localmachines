var actionTypes = require('../actionTypes');
var Cookies = require('js-cookie');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;


const initialState = toImmutable({
  locations: [],
  locationId: parseInt(Cookies.get('location'))
});

var LocationStore = new Nuclear.Store({
  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.SET_LOCATIONS, setLocations);
    this.on(actionTypes.SET_LOCATION_ID, setLocationId);
  }
});

function setLocations(state, { locations }) {
  return state.set('locations', locations);
}

function setLocationId(state, { id }) {
  return state.set('locationId', id);
}

export default LocationStore;

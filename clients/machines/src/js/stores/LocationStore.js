var actionTypes = require('../actionTypes');
var Cookies = require('js-cookie');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;

function queryLocation() {
  var query = window.location.search;
  if (query) {
    query = query.slice(1);
    var i = query.indexOf('location=');
    var location = query.slice(i + 'location='.length);
    var j = query.indexOf('&');
    if (j > 0) {
      location = location.slice(0, j);
    }
    return location;
  }
}

const initialState = toImmutable({
  locations: [],
  locationId: parseInt(queryLocation() || Cookies.get('location') || 1)
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

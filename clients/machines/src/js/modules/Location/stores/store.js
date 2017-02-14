import actionTypes from '../actionTypes';
import Cookies from 'js-cookie';
import Nuclear from 'nuclear-js';
var toImmutable = Nuclear.toImmutable;

const initialState = toImmutable({
  locations: [],
  locationId: parseInt(Cookies.get('location')) || 1
});

var LocationStore = new Nuclear.Store({
  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.SET_LOCATIONS, setLocations);
    this.on(actionTypes.SET_LOCATION_ID, setLocationId);
    this.on(actionTypes.SET_LOCATION_TERMS_URL, setLocationTermsUrl);
    this.on(actionTypes.SET_USER_LOCATIONS, setUserLocations);
  }
});

function setLocations(state, { locations }) {
  return state.set('locations', toImmutable(locations));
}

function setLocationId(state, { id }) {
  return state.set('locationId', id);
}

function setLocationTermsUrl(state, termsUrl) {
  return state.set('termsUrl', termsUrl);
}

function setUserLocations(state, userLocations) {
  return state.set('userLocations', toImmutable(userLocations));
}

export default LocationStore;

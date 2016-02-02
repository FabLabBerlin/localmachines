var $ = require('jquery');
var actionTypes = require('../actionTypes');
var reactor = require('../reactor');

function showGlobalLoader() {
  reactor.dispatch(actionTypes.SET_LOADING);
}

function hideGlobalLoader() {
  reactor.dispatch(actionTypes.UNSET_LOADING);
}

function loadAvailableLocations() {
  $.ajax({
    url: '/api/locations',
    success(locations) {
      reactor.dispatch(actionTypes.SET_LOCATIONS, { locations });
    },
    error(xhr, status, err) {
      toastr.error('Error loading locations');
      console.error(status, err);
    }
  });
}

export default {
  showGlobalLoader, hideGlobalLoader, loadAvailableLocations
};

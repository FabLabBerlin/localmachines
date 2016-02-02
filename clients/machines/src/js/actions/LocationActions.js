var $ = require('jquery');
var actionTypes = require('../actionTypes');
var reactor = require('../reactor');

var LocationActions = {

  loadLocations() {
    $.ajax({
      url: '/api/locations',
      success(locations) {
        reactor.dispatch(actionTypes.SET_LOCATIONS, { locations });
      },
      error(xhr, status, err) {
        toastr.error('Error loading locations');
      }
    });
  },

  setLocationId(id) {
    reactor.dispatch(actionTypes.SET_LOCATION_ID, { id });
  }

};





export default LocationActions;

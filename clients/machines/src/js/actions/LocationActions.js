var $ = require('jquery');
var actionTypes = require('../actionTypes');
var Cookies = require('js-cookie');
var reactor = require('../reactor');
var toastr = require('../toastr');

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
    console.log('LocationActions: setLocationId: ', id);
    Cookies.set('location', String(id));
    reactor.dispatch(actionTypes.SET_LOCATION_ID, { id });
  }

};





export default LocationActions;

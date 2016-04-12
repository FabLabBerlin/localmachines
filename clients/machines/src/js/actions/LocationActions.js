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

  loadUserLocations(userId) {
    var url = '/api/users/' + userId + '/locations';
    $.ajax({
      url: url,
      dataType: 'json',
      success(userLocations) {
        reactor.dispatch(actionTypes.SET_USER_LOCATIONS, userLocations);
      },
      error(xhr, status, err) {
        toastr.error('Error.  Please try again later.');
        console.error(url, status, err);
      }
    });
  },

  loadTermsUrl(locationId) {
    $.ajax({
      url: '/api/settings/terms_url?location=' + locationId,
      success(termsUrl) {
        reactor.dispatch(actionTypes.SET_LOCATION_TERMS_URL, termsUrl);
      },
      error(xhr, status, err) {
        toastr.error('Error loading terms');
      }
    });
  },

  addLocation({locationId, userId, router}) {
    var url = '/api/users/' + userId + '/locations/' + locationId;
    $.ajax({
      url: url,
      dataType: 'json',
      type: 'POST',
      success(data) {
        router.transitionTo('/machine');
      },
      error(xhr, status, err) {
        toastr.error('Error.  Please try again later.');
        console.error(url, status, err);
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

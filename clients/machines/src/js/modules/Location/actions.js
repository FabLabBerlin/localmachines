var $ = require('jquery');
var actionTypes = require('./actionTypes');
var Cookies = require('js-cookie');
var getters = require('../../getters');
var LocationGetters = require('./getters');
var reactor = require('../../reactor');
var toastr = require('../../toastr');

import {hashHistory} from 'react-router';


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
        console.log('loadUserLocations: xhr=', xhr);
        toastr.error('Error.  Please try again later.');
        console.log('loadUserLocations:', url, status, err);
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
        hashHistory.push('/machines');
      },
      error(xhr, status, err) {
        toastr.error('Error.  Please try again later.');
        console.error(url, status, err);
      }
    });
  },

  setLocationId(id) {
    reactor.dispatch(actionTypes.SET_LOCATION_ID, { id });
  },

  /*
   * Editing functions
   */

  saveEditedLocation() {
    var l = reactor.evaluateToJS(LocationGetters.getEditLocation);

    $.ajax({
      url: '/api/locations/',
      dataType: 'json',
      type: 'PUT'
    })
    .done(() => {
      toastr.info('Successfully updated location.');
      LocationActions.loadLocations();
    })
    .fail(() => {
      toastr.error('Error updating location.  Please try again later.');
    });
  },

  setEditLocation(location) {
    reactor.dispatch(actionTypes.SET_EDIT_LOCATION, location);
  }

};





export default LocationActions;

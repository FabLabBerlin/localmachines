import $ from 'jquery';
import actionTypes from './actionTypes';
import Cookies from 'js-cookie';
import getters from '../../getters';
import LocationGetters from './getters';
import reactor from '../../reactor';
import toastr from '../../toastr';

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

  addEditLocation() {
    console.log('addEditLocation()');
    $.ajax({
      url: '/api/locations',
      dataType: 'json',
      type: 'POST',
      data: {
        title: 'Untitled'
      }
    })
    .done(() => {
      toastr.info('Successfully added location.');
      LocationActions.loadLocations();
    })
    .fail(() => {
      toastr.error('Error adding location.  Please try again later.');
    });
  },

  saveEditedLocation() {
    var l = reactor.evaluateToJS(LocationGetters.getEditLocation);

    $.ajax({
      url: '/api/locations/' + l.Id,
      dataType: 'json',
      type: 'PUT',
      contentType: 'application/json; charset=utf-8',
      data: JSON.stringify(l)
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

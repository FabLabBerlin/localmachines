var _ = require('lodash');


const getLocations = [
  ['locationStore'],
  (locationStore) => {
    return locationStore.get('locations');
  }
];

const getLocation = [
  getLocations,
  ['locationStore'],
  (locations, locationStore) => {
    var location;
    _.each(locations, (l) => {
      if (l.Id === locationStore.get('locationId')) {
        location = l;
      }
    });
    return location;
  }
];

const getLocationId = [
  ['locationStore'],
  (locationStore) => {
    return locationStore.get('locationId');
  }
];

const getLocationTermsUrl = [
  ['locationStore'],
  (locationStore) => {
    return locationStore.get('termsUrl');
  }
];

const getUserLocations = [
  ['locationStore'],
  (locationStore) => {
    return locationStore.get('userLocations');
  }
];

const getUserLocation = [
  getLocationId,
  getUserLocations,
  ['locationStore'],
  (locationId, userLocations, locationStore) => {
    if (userLocations) {
      var userLocation = userLocations.find((ul) => {
        return ul.get('LocationId') === locationId;
      });
      return userLocation;
    }
  }
];

const getIsStaff = [
  getUserLocation,
  (userLocation) => {
    if (userLocation) {
      var role = userLocation.get('UserRole');
      return role === 'staff' || role === 'admin' || role === 'superadmin';
    }
  }
];

const getIsAdmin = [
  getUserLocation,
  (userLocation) => {
    if (userLocation) {
      var role = userLocation.get('UserRole');
      return role === 'admin' || role === 'superadmin';
    }
  }
];

export default {
  getLocation, getLocationId, getLocations, getLocationTermsUrl,
  getIsStaff, getIsAdmin
};
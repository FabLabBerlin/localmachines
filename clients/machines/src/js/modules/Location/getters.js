import _ from 'lodash';


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
    if (locations) {
      return locations.find((l) => {
        return l.get('Id') === locationStore.get('locationId');
      });
    }
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

const getEditLocation = [
  ['locationEditStore'],
  (locationEditStore) => {
    return locationEditStore;
  }
];

export default {
  getLocation, getLocationId, getLocations, getLocationTermsUrl,
  getIsStaff, getIsAdmin, getUserLocation,
  getEditLocation
};

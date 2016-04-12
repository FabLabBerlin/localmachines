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

export default {
  getLocation, getLocationId, getLocations, getLocationTermsUrl
};

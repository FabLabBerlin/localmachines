const getAllMemberships = [
  ['membershipsStore'],
  (membershipsStore) => {
    return membershipsStore.get('allMemberships');
  }
];

const getShowArchived = [
  ['membershipsStore'],
  (membershipsStore) => {
    return membershipsStore.get('showArchived');
  }
];

export default {
  getAllMemberships,
  getShowArchived
};

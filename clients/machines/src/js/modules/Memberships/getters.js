const getAllMemberships = [
  ['membershipsStore'],
  (membershipsStore) => {
    return membershipsStore.get('allMemberships');
  }
];

export default {
  getAllMemberships
};

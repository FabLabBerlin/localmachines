const getMemberships = [
  ['usersStore'],
  (usersStore) => {
    return usersStore.get('memberships');
  }
];

const getUserMemberships = [
  ['usersStore'],
  (usersStore) => {
    return usersStore.get('userMemberships');
  }
];

const getUsers = [
  ['usersStore'],
  (usersStore) => {
    return usersStore.get('users');
  }
];

export default {
  getMemberships,
  getUserMemberships,
  getUsers
};

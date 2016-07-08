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
  getUserMemberships,
  getUsers
};

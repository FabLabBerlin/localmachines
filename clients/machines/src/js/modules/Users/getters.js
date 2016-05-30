const getUsers = [
  ['usersStore'],
  (usersStore) => {
    return usersStore.get('users');
  }
];

export default {
  getUsers
};

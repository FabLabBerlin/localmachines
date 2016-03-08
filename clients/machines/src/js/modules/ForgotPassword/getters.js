const getKey = [
  ['forgotPasswordStore'],
  (forgotPasswordStore) => {
    return forgotPasswordStore.get('key');
  }
];

const getPhone = [
  ['forgotPasswordStore'],
  (forgotPasswordStore) => {
    return forgotPasswordStore.get('phone');
  }
];

export default {
  getKey,
  getPhone
};

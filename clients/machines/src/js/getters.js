const getIsLogged = [
  ['loginStore'],
  (loginStore) => {
    return loginStore.get('isLogged');
  }
];

const getUid = [
  ['loginStore'],
  (loginStore) => {
    return loginStore.get('uid');
  }
];

const getFirstTry = [
  ['loginStore'],
  (loginStore) => {
    return loginStore.get('firstTry');
  }
];

const getLoginSuccess = [
  ['loginStore'],
  (loginStore) => {
    return loginStore.get('loginSuccess');
  }
];

const getUserInfo = [
  ['machineStore'],
  (machineStore) => {
    return machineStore.get('userInfo');
  }
];

const getActivationInfo = [
  ['machineStore'],
  (machineStore) => {
    return machineStore.get('activationInfo');
  }
];

const getMachineInfo = [
  ['machineStore'],
  (machineStore) => {
    return machineStore.get('machineInfo') || [];
  }
];

export default {
  getIsLogged, getUid, getFirstTry, getLoginSuccess,
  getUserInfo, getActivationInfo, getMachineInfo
};

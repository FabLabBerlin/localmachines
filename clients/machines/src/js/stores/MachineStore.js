var $ = require('jquery');
var _ = require('lodash');
var actionTypes = require('../actionTypes');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;


const initialState = toImmutable({
  userInfo: {},
  activationInfo: [],
  machineInfo: []
});

/*
 * TODO: refactoring some comments and function
 * Store the data
 * summary:
 * state (:34)
 * postAPI template function (:46)
 * getAPI template function (:69)
 * Call order (callback are define below or inside the apicall):
 *  - userInfo (:128)
 *  - machineInfo (:151)
 *  - getActivationInfo (:167)
 *  - postActivationInfo (:190)
 *  - putActivation (:215)
 *  - postSwitchMachine (:238)
 * LoginError (:273)
 * utils:
 *  - formatActivation (:285)
 *  - nameInAllActivation (:301)
 *  - nameInOneActivation (:311)
 *  - getter (:324)
 *  - putLoginState (:343)
 *  - onChange (:353)
 */
var MachineStore = new Nuclear.Store({
  /*
   * State of MachineStore
   * Information needed by the components
   */

  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.MACHINE_STORE_CLEAR_STATE, clearState);
    this.on(actionTypes.SET_ACTIVATION_INFO, setActivationInfo);
    this.on(actionTypes.SET_USER_INFO, setUserInfo);
    this.on(actionTypes.SET_MACHINE_INFO, setMachineInfo);
  }

});

/*
 * Clean State before login out
 */
function clearState(state) {
  return initialState;
}

function setActivationInfo(state, { activationInfo }) {
  return state.set('activationInfo', activationInfo);
}

function setMachineInfo(state, { machineInfo }) {
  return state.set('machineInfo', machineInfo);
}

function setUserInfo(state, { userInfo }) {
  return state.set('userInfo', userInfo);
}

export default MachineStore;

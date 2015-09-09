var $ = require('jquery');
var _ = require('lodash');
var actionTypes = require('../actionTypes');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;


const initialState = toImmutable({
  activationInfo: [],
  machineInfo: [],
  machineUsers: {},
  loading: false
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
    this.on(actionTypes.REGISTER_MACHINE_USER, registerMachineUser);
    this.on(actionTypes.SET_ACTIVATION_INFO, setActivationInfo);
    this.on(actionTypes.SET_MACHINE_INFO, setMachineInfo);
    this.on(actionTypes.SET_UNDER_MAINTENANCE, setUnderMaintenance);
    this.on(actionTypes.SET_LOADING, setLoading);
    this.on(actionTypes.UNSET_LOADING, unsetLoading);
  }

});

/*
 * Clean State before login out
 */
function clearState(state) {
  return initialState;
}

function registerMachineUser(state, { userData }) {
  return state.set('machineUsers', state.get('machineUsers').set(userData.UserId, userData));
}

function setActivationInfo(state, { activationInfo }) {
  return state.set('activationInfo', activationInfo);
}

function setMachineInfo(state, { machineInfo }) {
  const machinesById = toImmutable(machineInfo)
    .toMap()
    .mapKeys((k, v) => v.get('Id'));
  return state.set('machinesById', machinesById);
}

function setUnderMaintenance(state, { mid, onOrOff }) {
  var m = state.get('machinesById').get(mid)
                                   .set('UnderMaintenance', onOrOff === 'on');
  return state.set('machinesById', state.get('machinesById')
                                        .set(mid, m));
}

function setLoading(state) {
  return state.set('loading', true);
}

function unsetLoading(state) {
  return state.set('loading', false);
}

export default MachineStore;

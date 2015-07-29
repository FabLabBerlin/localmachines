import MachineStore from '../stores/MachineStore';

var MachineActions = {

  fetchData(uid) {
    MachineStore.apiGetUserInfoLogin(uid);
  },

  /*
   * To end an activation
   * @aid: id of the activation you want to shut down
   */
  endActivation(aid) {
    MachineStore.apiPutActivation(aid);
  },

  /*
   * To start an activation
   * @mid: id of the machine you want to activate
   */
  startActivation(mid) {
    MachineStore.apiPostActivation(mid);
  },

  /*
   * When an admin want to force on a machine
   */
  adminTurnOffMachine(mid, aid) {
    MachineStore.apiPostSwitchMachine(mid, 'off', aid);
  },

  /*
   * When an admin want to force off a machine
   */
  adminTurnOnMachine(mid) {
    MachineStore.apiPostSwitchMachine(mid, 'on');
  },

  /*
   * Clear store state while logout
   */
  clearState() {
    MachineStore.clearState();
  },

  /*
   * To continue to refresh the view each seconds
   */
  pollActivations() {
    MachineStore.apiGetActivationActive();
  }

};

module.exports = MachineActions;

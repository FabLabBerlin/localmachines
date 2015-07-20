import MachineStore from '../stores/MachineStore';

var MachineActions = {

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
  adminTurnOnMachine(mid) {
    MachineStore.apiPostSwitchMachine(mid, 'on', aid);
  },

  /*
   * When an admin want to force off a machine
   */
  adminTurnOffMachine(mid) {
    MachineStore.apiPostSwitchMachine(mid, 'off');
  },

  /*
   * To continue to refresh the view each seconds
   */
  pollActivations() {
    MachineStore.apiGetActivationActive();
  }

}

module.exports = MachineActions;

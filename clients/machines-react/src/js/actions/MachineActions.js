import MachineStore from '../stores/MachineStore';

var MachineActions = {

  /*
   * To end an activation
   * @mid: machine id you want to turn off
   * @aid: id of the activation you want to shut down
   */
  endActivation(mid, aid) {
    MachineStore.apiPostSwitchMachine(mid, 'off', aid);
  },

  /*
   * To start an activation
   * @mid: id of the machine you want to activate
   */
  startActivation(mid) {
    MachineStore.apiPostSwitchMachine(mid, 'on');
  },

  /*
   * To continue to refresh the view each seconds
   */
  pollActivations() {
    MachineStore.apiGetActivationActive();
    setTimeout(function() {
      this.pollActivations();
    }.bind(this), 1000);
  }

}

module.exports = MachineActions;

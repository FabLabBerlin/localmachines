import MachineStore from '../stores/MachineStore';

var MachineActions = {

  /*
   * To end an activation
   * @aid: id of the activation you want to shut down
   */
  endActivation(aid) {
    MachineStore.putActivation(aid);
  },

  /*
   * To start an activation
   * @mid: id of the machine you want to activate
   */
  startActivation(mid) {
    MachineStore.postActivation(mid);
  },

  /*
   * To continue to refresh the view each seconds
   */
  pollActivations() {
    MachineStore.getActivationActive();
    setTimeout(function() {
      this.pollActivations();
    }.bind(this), 1000);
  }

}

module.exports = MachineActions;

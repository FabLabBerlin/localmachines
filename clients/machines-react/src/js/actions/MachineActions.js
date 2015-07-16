import MachineStore from '../stores/MachineStore';

var MachineActions = {

  endActivation(aid) {
    MachineStore.endActivation(aid);
  },

  startActivation(mid) {
    MachineStore.postActivation(mid);
  }

}

module.exports = MachineActions;

import MachineStore from '../stores/MachineStore';

var MachineActions = {

  endActivation(aid) {
    MachineStore.putActivation(aid);
  },

  startActivation(mid) {
    MachineStore.postActivation(mid);
  }

}

module.exports = MachineActions;

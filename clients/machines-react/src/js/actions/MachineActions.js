import MachineStore from '../stores/MachineStore';

var MachineActions = {

  endActivation(aid) {
    MachineStore.endActivation(aid);
  },

  startActivation(mid) {
    MachineStore.startActivation(mid);
  }

}

module.exports = MachineActions;

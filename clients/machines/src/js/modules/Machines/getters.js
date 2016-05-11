var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;


const getActivations = [
  ['machineStore'],
  (machineStore) => {
    return machineStore.get('activations');
  }
];

const getMachinesById = [
  ['machineStore'],
  (machineStore) => {
    return machineStore.get('machinesById') || toImmutable({});
  }
];

const getMachines = [
  getMachinesById,
  (machinesById) => {
    return machinesById.sortBy(m => m.Name);
  }
];

const getMachineUsers = [
  ['machineStore'],
  (machineStore) => {
    return machineStore.get('machineUsers') || {};
  }
];

export default {
  getActivations, getMachinesById, getMachines, getMachineUsers
};

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

export default {
  getActivations, getMachinesById
};

const getActivations = [
  ['machineStore'],
  (machineStore) => {
    return machineStore.get('activations');
  }
];

export default {
  getActivations
};

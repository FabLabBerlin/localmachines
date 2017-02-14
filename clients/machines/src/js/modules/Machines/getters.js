import LocationGetters from '../Location/getters';
import Nuclear from 'nuclear-js';
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

const getMyMachines = [
  ['machineStore'],
  getActivations,
  getMachinesById,
  ['loginStore'],
  LocationGetters.getLocationId,
  (machineStore, activations, machinesById, loginStore, locationId) => {
    const uid = loginStore.get('uid');
    if (activations && machinesById) {
      return activations
        .filter(a => a.get('UserId') === uid && a.get('LocationId') === locationId)
        .map(a => machinesById.get(a.get('MachineId')));
    }
  }
];

const getNewMachineImages = [
  ['machineStore'],
  (machineStore) => {
    return machineStore.get('newMachineImages');
  }
];

export default {
  getActivations,
  getMachinesById,
  getMachines,
  getMachineUsers,
  getMyMachines,
  getNewMachineImages
};

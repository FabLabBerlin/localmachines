var $ = require('jquery');
var _ = require('lodash');
var actionTypes = require('../actionTypes');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;


const initialState = toImmutable({
  activations: [],
  locationId: 1,
  machines: [],
  machineUsers: {}
});


var MachineStore = new Nuclear.Store({

  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.MACHINE_STORE_CLEAR_STATE, clearState);
    this.on(actionTypes.REGISTER_MACHINE_USER, registerMachineUser);
    this.on(actionTypes.SET_ACTIVATIONS, setActivations);
    this.on(actionTypes.SET_LOCATION_ID, setLocationId);
    this.on(actionTypes.SET_LOCATIONS, setLocations);
    this.on(actionTypes.SET_MACHINES, setMachines);
    this.on(actionTypes.SET_UNDER_MAINTENANCE, setUnderMaintenance);
  }

});

function clearState(state) {
  return initialState;
}

function registerMachineUser(state, { userData }) {
  return state.set('machineUsers', state.get('machineUsers').set(parseInt(userData.UserId, 10), userData));
}

function setActivations(state, { activations }) {
  return state.set('activations', activations);
}

function setLocationId(state, { id }) {
  return state.set('locationId', id);
}

function setLocations(state, { locations }) {
  return state.set('locations', locations);
}

function setMachines(state, { machines }) {
  const machinesById = toImmutable(machines || [])
    .toMap()
    .mapKeys((k, v) => v.get('Id'));
  return state.set('machinesById', machinesById);
}

function setUnderMaintenance(state, { mid, onOrOff }) {
  var m = state.get('machinesById').get(mid)
                                   .set('UnderMaintenance', onOrOff === 'on');
  return state.set('machinesById', state.get('machinesById')
                                        .set(mid, m));
}

export default MachineStore;

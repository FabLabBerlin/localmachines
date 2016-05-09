var $ = require('jquery');
var _ = require('lodash');
var actionTypes = require('../actionTypes');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;


const initialState = toImmutable({
  activations: [],
  machines: [],
  machineUsers: {}
});


var MachineStore = new Nuclear.Store({

  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.MACHINE_STORE_CLEAR_STATE, clearState);
    this.on(actionTypes.REGISTER_MACHINE_USERS, registerMachineUsers);
    this.on(actionTypes.SET_ACTIVATIONS, setActivations);
    this.on(actionTypes.SET_MACHINES, setMachines);
    this.on(actionTypes.SET_UNDER_MAINTENANCE, setUnderMaintenance);
  }

});

function clearState(state) {
  return initialState;
}

function registerMachineUsers(state, users) {
  _.each(users, function(user) {
    state = state.set('machineUsers', state.get('machineUsers').
      set(parseInt(user.UserId, 10), user));
  });
  return state;
}

function setActivations(state, { activations }) {
  return state.set('activations', toImmutable(activations));
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

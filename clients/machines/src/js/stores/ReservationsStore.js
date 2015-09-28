var actionTypes = require('../actionTypes');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;


const initialState = toImmutable({
  reservations: undefined,
  create: null
});

var ReservationsStore = new Nuclear.Store({
  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.SET_RESERVATIONS, setReservations);
    this.on(actionTypes.CREATE_EMPTY, createEmpty);
    this.on(actionTypes.CREATE_SET_MACHINE, createSetMachine);
    this.on(actionTypes.CREATE_SET_DATE, createSetDate);
    this.on(actionTypes.CREATE_SET_TIME_FROM, createSetTimeFrom);
    this.on(actionTypes.CREATE_SET_TIME_TO, createSetTimeTo);
  }
});

function setReservations(state, { reservations }) {
  return state.set('reservations', toImmutable(reservations));
}

function createEmpty(state) {
  return state.set('create', toImmutable({}));
}

function createSetMachine(state, { mid }) {
  return state.setIn(['create', 'machineId'], mid);
}

function createSetDate(state, { date }) {
  return state.setIn(['create', 'date'], date);
}

function createSetTimeFrom(state, { timeFrom }) {
  return state.setIn(['create', 'timeFrom'], timeFrom);
}

function createSetTimeTo(state, { timeTo }) {
  return state.setIn(['create', 'timeTo'], timeTo);
}

export default ReservationsStore;

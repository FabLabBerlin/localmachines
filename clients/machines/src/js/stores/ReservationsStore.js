import actionTypes from '../actionTypes';
import Nuclear from 'nuclear-js';
var toImmutable = Nuclear.toImmutable;


/* Create Steps */
const STEP_SET_MACHINE = 1;
const STEP_SET_DATE = 2;
const STEP_SET_TIME = 3;
const STEP_SUCCESS = 4;
const STEP_ERROR = 5;

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
    this.on(actionTypes.CREATE_SET_TIMES, createSetTimes);
    this.on(actionTypes.CREATE_DONE, createDone);
    this.on(actionTypes.CREATE_SET_STEP, createSetStep);
  }
});

function setReservations(state, { reservations }) {
  return state.set('reservations', toImmutable(reservations));
}

function createEmpty(state) {
  return state.set('create', toImmutable({
    step: STEP_SET_MACHINE
  }));
}

function createSetMachine(state, { mid }) {
  return state.setIn(['create', 'machineId'], mid);
}

function createSetDate(state, { date }) {
  state = state.setIn(['create', 'date'], date);
  state = state.setIn(['create', 'times'], toImmutable([]));
  return state;
}

function createSetTimes(state, { times }) {
  return state.setIn(['create', 'times'], toImmutable(times));
}

function createDone(state) {
  return state.set('create', null);
}

function createSetStep(state, step) {
  return state.setIn(['create', 'step'], step);
}

export default ReservationsStore;

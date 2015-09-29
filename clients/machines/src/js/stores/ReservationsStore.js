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
    this.on(actionTypes.CREATE_SET_TIMES, createSetTimes);
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
  state = state.setIn(['create', 'date'], date);
  state = state.setIn(['create', 'times'], possibleTimes(state));
  return state;
}

function possibleTimes(state) {
  var date = state.get('create').get('date');
  var times = [];
  var startHour = 10;
  var endHour = 19;
  if (date.isoWeekday() === 6) {
    startHour = 12;
    endTime = 18;
  }
  if (date.isoWeekday() !== 7) {
    for (var tStart = date.clone().hours(startHour); tStart.hours() < endHour; tStart.add(30, 'm')) {
      times.push({
        start: tStart.clone(),
        end: tStart.clone().add(30, 'm'),
        selected: false
      });
    }
  }
  return toImmutable(times);
}

function createSetTimes(state, { times }) {
  return state.setIn(['create', 'times'], toImmutable(times));
}

export default ReservationsStore;

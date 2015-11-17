var actionTypes = require('../actionTypes');
var ApiActions = require('./ApiActions');
var getters = require('../getters');
var moment = require('moment');
var reactor = require('../reactor');
var toastr = require('toastr');

const STEP_SET_MACHINE = 1;
const STEP_SET_DATE = 2;
const STEP_SET_TIME = 3;
const STEP_SUCCESS = 4;
const STEP_ERROR = 5;

var ReservationActions = {
  STEP_SET_MACHINE: STEP_SET_MACHINE,
  STEP_SET_DATE: STEP_SET_DATE,
  STEP_SET_TIME: STEP_SET_TIME,
  STEP_SUCCESS: STEP_SUCCESS,
  STEP_ERROR: STEP_ERROR,

  load() {
    ApiActions.getCall('/api/reservations', function(reservations) {
      _.each(reservations, function(reservation) {
        ApiActions.getCall('/api/users/' + reservation.UserId + '/name', function(userData) {
          reactor.dispatch(actionTypes.REGISTER_MACHINE_USER, { userData });
        });
      });
      reactor.dispatch(actionTypes.SET_RESERVATIONS, { reservations });
    });
  },

  createEmpty() {
    reactor.dispatch(actionTypes.CREATE_EMPTY);
  },

  createSetMachine({ mid }) {
    reactor.dispatch(actionTypes.CREATE_SET_MACHINE, { mid });
  },

  createSetDate({ date }) {
    reactor.dispatch(actionTypes.CREATE_SET_DATE, { date });
    reactor.dispatch(actionTypes.CREATE_SET_STEP, STEP_SET_TIME);
  },

  createSetTimes({ times }) {
    reactor.dispatch(actionTypes.CREATE_SET_TIMES, { times });
  },

  createToggleStartTime({ startTime }) {
    reactor.dispatch(actionTypes.CREATE_TOGGLE_START_TIME);
  },

  previousStep() {
    var step = reactor.evaluateToJS(getters.getNewReservation).step;
    if (step > STEP_SET_MACHINE && step !== STEP_ERROR) {
      step--;
    }
    reactor.dispatch(actionTypes.CREATE_SET_STEP, step);
  },

  nextStep() {
    var newReservation = reactor.evaluateToJS(getters.getNewReservation);
    var step = newReservation.step;
    if (step === STEP_SET_DATE && !newReservation.date) {
      toastr.error('Please select a date');
    } else if (step < STEP_SUCCESS && step !== STEP_ERROR) {
      step++;
      reactor.dispatch(actionTypes.CREATE_SET_STEP, step);
    }
  },

  createSubmit() {
    const times = reactor.evaluateToJS(getters.getNewReservationTimes);
    const reservation = reactor.evaluateToJS(getters.getNewReservation);
    const uid = reactor.evaluateToJS(getters.getUid);
    var selectedTimes = _.filter(times, function(t) {
      return t.selected;
    });
    var data = {
      MachineId: reservation.machineId,
      UserId: uid,
      TimeStart: reactor.evaluateToJS(getters.getNewReservationFrom),
      TimeEnd: reactor.evaluateToJS(getters.getNewReservationTo),
      Created: new Date()
    };
    $.ajax({
      url: '/api/reservations',
      contentType: 'application/json; charset=utf-8',
      dataType: 'json',
      type: 'POST',
      cache: false,
      data: JSON.stringify(data),
      success: function() {
        reactor.dispatch(actionTypes.CREATE_SET_STEP, STEP_SUCCESS);
        ReservationActions.load();
      },
      error: function() {
        toastr.error('Error submitting reservation. Please try again later.');
      }
    });
  },

  createDone() {
    reactor.dispatch(actionTypes.CREATE_DONE);
  },

  deleteReservation(reservationId) {
    reactor.dispatch(actionTypes.DELETE_RESERVATION_START);

    $.ajax({
      url: '/api/reservations/' + reservationId,
      type: 'DELETE',
      cache: false,
      success: function() {
        reactor.dispatch(actionTypes.DELETE_RESERVATION_SUCCESS);
        toastr.success('Successfuly deleted reservation');
        ReservationActions.load();
      },
      error: function() {
        reactor.dispatch(actionTypes.DELETE_RESERVATION_FAIL);
        toastr.error('Error deleting reservation');
      }
    });
  }
};

export default ReservationActions;

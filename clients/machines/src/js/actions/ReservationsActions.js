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
      reactor.dispatch(actionTypes.SET_RESERVATIONS, { reservations });
    });
  },

  createEmpty() {
    reactor.dispatch(actionTypes.CREATE_EMPTY);
  },

  createSetMachine({ mid }) {
    reactor.dispatch(actionTypes.CREATE_SET_MACHINE, { mid });
    reactor.dispatch(actionTypes.CREATE_SET_STEP, STEP_SET_DATE);
  },

  createSetDate({ date }) {
    date = moment(date, 'YYYY-MM-DD');
    if (!date.isValid()) {
      toastr.error('Please enter date in the format YYYY-MM-DD');
    } else if (date.isBefore(moment())) {
      toastr.error('Please enter date from the future');
    } else if (date.isAfter(moment().add(2, 'months'))) {
      toastr.error('Please enter date within the next 2 months');
    } else if (date.isoWeekday() === 7) {
      toastr.error('Please enter a weekday');
    } else {
      reactor.dispatch(actionTypes.CREATE_SET_DATE, { date });
      reactor.dispatch(actionTypes.CREATE_SET_STEP, STEP_SET_TIME);
    }
  },

  createSetTimes({ times }) {
    reactor.dispatch(actionTypes.CREATE_SET_TIMES, { times });
  },

  createToggleStartTime({ startTime }) {
    reactor.dispatch(actionTypes.CREATE_TOGGLE_START_TIME);
  },

  createSubmit() {
    const times = reactor.evaluateToJS(getters.getNewReservationTimes);
    const reservation = reactor.evaluateToJS(getters.getNewReservation);
    const uid = reactor.evaluateToJS(getters.getUid);
    if (!this.isRange(times)) {
      toastr.error('Please select a time range');
      return;
    }
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

  isRange(times) {
    var lastIsSelected;
    var rangesFound = 0;
    _.each(times, function(t) {
      if (t.selected && !lastIsSelected) {
        rangesFound++;
      }
      lastIsSelected = t.selected;
    });
    return rangesFound === 1;
  }
};

export default ReservationActions;

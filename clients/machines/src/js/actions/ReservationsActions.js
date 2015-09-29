var actionTypes = require('../actionTypes');
var ApiActions = require('./ApiActions');
var getters = require('../getters');
var moment = require('moment');
var reactor = require('../reactor');
var toastr = require('toastr');


var ReservationActions = {

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
  },

  createSetDate({ date }) {
    date = moment(date, 'YYYY-MM-DD');
    if (!date.isValid()) {
      toastr.error('Please enter date in the format YYYY-MM-DD');
    } else if (date.isBefore(moment())) {
      toastr.error('Please enter date from the future');
    } else if (date.isoWeekday() === 7) {
      toastr.error('Please enter a weekday');
    } else {
      reactor.dispatch(actionTypes.CREATE_SET_DATE, { date });
    }
  },

  createSetTimes({ times }) {
    reactor.dispatch(actionTypes.CREATE_SET_TIMES, { times });
    const reservation = reactor.evaluateToJS(getters.getNewReservation);
    const uid = reactor.evaluateToJS(getters.getUid);
    console.log('createSetTimes...reservation:', reservation);
    $.ajax({
      url: '/api/reservations',
      contentType: "application/json; charset=utf-8",
      dataType: 'json',
      type: 'POST',
      data: JSON.stringify({
        MachineId: reservation.machineId,
        UserId: uid,
        TimeStart: times[0].start.toDate(),
        TimeEnd: times[times.length - 1].start.toDate(),
        Created: new Date()
      })
    });
  },

  createToggleStartTime({ startTime }) {
    reactor.dispatch(actionTypes.CREATE_TOGGLE_START_TIME);
  }
};

export default ReservationActions;

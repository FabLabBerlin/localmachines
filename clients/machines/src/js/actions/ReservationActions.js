import _ from 'lodash';
var $ = require('jquery');
import actionTypes from '../actionTypes';
import Cookies from 'js-cookie';
import getters from '../getters';
import Machines from '../modules/Machines';
import moment from 'moment';
import reactor from '../reactor';
import toastr from 'toastr';

import helpers from '../components/UserProfile/helpers';


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

  // load reservations and corresponding user names.
  //
  // But only load user names for reservations that are currently active.
  load() {
    const lid = Cookies.get('location');
    $.ajax({
      url: '/api/reservations?location=' + lid
    })
    .done(reservations => {
      var t = moment().unix();
      var userIds = [];

      _.each(reservations, function(reservation) {
        var twoDays = 2 * 86400;
        var timeStart = moment(reservation.TimeStart).unix() - twoDays;
        var timeEnd = helpers.timeEnd(reservation).unix() + twoDays;
        if (timeStart <= t && t <= timeEnd) {
          userIds.push(reservation.UserId);
        }
      });

      $.ajax({
        url: '/api/users/names?uids=' + userIds.join(','),
        dataType: 'json',
        type: 'GET',
        success(data) {
          reactor.dispatch(Machines.actionTypes.REGISTER_MACHINE_USERS, data.Users);
        },
        error() {
            console.log('Error loading names');
        }
      });

      reactor.dispatch(actionTypes.SET_RESERVATIONS, { reservations });
    });
  },

  newReservation: {
    create() {
      reactor.dispatch(actionTypes.CREATE_EMPTY);
    },

    setMachine({ mid }) {
      reactor.dispatch(actionTypes.CREATE_SET_MACHINE, { mid });
    },

    setDate({ date }) {
      reactor.dispatch(actionTypes.CREATE_SET_DATE, { date });
      reactor.dispatch(actionTypes.CREATE_SET_STEP, STEP_SET_TIME);
    },

    setTimes({ times }) {
      reactor.dispatch(actionTypes.CREATE_SET_TIMES, { times });
    },

    toggleStartTime({ startTime }) {
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

    submit() {
      const lid = parseInt(Cookies.get('location'));
      const times = reactor.evaluateToJS(getters.getNewReservationTimes);
      const reservation = reactor.evaluateToJS(getters.getNewReservation);
      const uid = reactor.evaluateToJS(getters.getUid);
      var selectedTimes = _.filter(times, function(t) {
        return t.selected;
      });
      var data = {
        LocationId: lid,
        MachineId: reservation.machineId,
        UserId: uid,
        TimeStart: reactor.evaluateToJS(getters.getNewReservationFrom),
        Quantity: reactor.evaluateToJS(getters.getNewReservationQuantity),
        Created: new Date()
      };
      $.ajax({
        url: '/api/reservations?location=' + lid,
        contentType: 'application/json; charset=utf-8',
        dataType: 'json',
        type: 'POST',
        cache: false,
        data: JSON.stringify(data),
        success() {
          reactor.dispatch(actionTypes.CREATE_SET_STEP, STEP_SUCCESS);
          ReservationActions.load();
        },
        error() {
          toastr.error('Error submitting reservation. Please try again later.');
        }
      });
    },

    done() {
      reactor.dispatch(actionTypes.CREATE_DONE);
    }
  },

  cancelReservation(reservationId) {
    const reservations = reactor.evaluateToJS(getters.getReservations);
    var reservation;
    _.each(reservations, (r) => {
      if (r.Id === reservationId) {
        reservation = r;
      }
    });

    if (!reservation) {
      return;
    }

    reactor.dispatch(actionTypes.CANCEL_RESERVATION_START);

    reservation.Cancelled = true;

    $.ajax({
      headers: {'Content-Type': 'application/json'},
      url: '/api/reservations/' + reservationId + '/cancel',
      dataType: 'json',
      type: 'POST',
      success() {
        reactor.dispatch(actionTypes.CANCEL_RESERVATION_SUCCESS);
        toastr.success('Successfuly cancelled reservation');
        ReservationActions.load();
      },
      error() {
        reactor.dispatch(actionTypes.CANCEL_RESERVATION_FAIL);
        toastr.error('Error cancelling reservation');
      }
    });
  }
};

export default ReservationActions;

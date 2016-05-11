jest.dontMock('../../actionTypes');
jest.dontMock('../../components/Reservations/helpers');
jest.dontMock('../../getters');
jest.dontMock('lodash');
jest.dontMock('../../modules/Machines/stores/store.js');
jest.dontMock('../../modules/Machines');
jest.dontMock('moment');
jest.dontMock('nuclear-js');
jest.dontMock('../../reactor');
jest.dontMock('../ReservationRulesStore.js');
jest.dontMock('../ReservationsStore.js');
jest.mock('jquery');
jest.mock('toastr');

var _ = require('lodash');
var actionTypes = require('../../actionTypes');
var getters = require('../../getters');
var Machines = require('../../modules/Machines');
var MachineStore = require('../../modules/Machines/stores/store.js');
var moment = require('moment');
var ReservationRulesStore = require('../ReservationRulesStore');
var ReservationsStore = require('../ReservationsStore');
var reactor = require('../../reactor');


function getMachines() {
  return {
    1: {
      Id: 1,
      Name: 'Printer123'
    }
  };
}

function reservationRules() {
  return [
    {  
      'MachineId': 0,
      'Created': '2015-10-05T14:28:51+02:00',
      'Available': true,
      'Name': 'Opening Hours Mo - Fr',
      'Tuesday': true,
      'DateEnd': '2015-12-31',
      'TimeEnd': '19:00',
      'Unavailable': false,
      'Monday': true,
      'Sunday': false,
      'DateStart': '2015-01-01',
      'TimeStart': '10:00',
      'Wednesday': true,
      'TimeZone': '',
      'Saturday': false,
      'Thursday': true,
      'Id': 100,
      'Friday': true
    },
    {  
      'MachineId': 0,
      'Created': '2015-10-05T14:28:51+02:00',
      'Available': true,
      'Name': 'Opening Hours Sa',
      'Tuesday': false,
      'DateEnd': '2015-12-31',
      'TimeEnd': '18:00',
      'Unavailable': false,
      'Monday': false,
      'Sunday': false,
      'DateStart': '2015-01-01',
      'TimeStart': '12:00',
      'Wednesday': false,
      'TimeZone': '',
      'Saturday': true,
      'Thursday': false,
      'Id': 101,
      'Friday': false
    }
  ];
}

describe('ReservationsStore', function() {
  reactor.registerStores({
    machineStore: MachineStore,
    reservationRulesStore: ReservationRulesStore,
    reservationsStore: ReservationsStore
  });

  describe('CREATE_SET_DATE', function() {
    it('sets the date and possible times at that day', function() {
      reactor.dispatch(actionTypes.SET_RESERVATION_RULES, reservationRules());
      reactor.dispatch(Machines.actionTypes.SET_MACHINES, {
        machines: getMachines()
      });
      reactor.dispatch(actionTypes.SET_RESERVATIONS, {
        reservations: []
      });
      reactor.dispatch(actionTypes.CREATE_EMPTY);
      reactor.dispatch(actionTypes.CREATE_SET_DATE, {
        date: moment('2015-10-16')
      });
      var possibleTimes = reactor.evaluateToJS(getters.getNewReservationTimes);
      expect(possibleTimes.length).toEqual(48);
      expect(_.first(possibleTimes).start.format('HH:mm')).toEqual('00:00');
      expect(_.first(possibleTimes).end.format('HH:mm')).toEqual('00:30');
      expect(_.last(possibleTimes).start.format('HH:mm')).toEqual('23:30');
      expect(_.last(possibleTimes).end.format('HH:mm')).toEqual('00:00');
      _.each(possibleTimes, (t, i) => {
        if (i >= 20 && i < 38) {
          expect(t.availableMachineIds).toEqual([1]);
        } else {
          expect(t.availableMachineIds).toEqual([]);
        }
      });
    });
  });
});

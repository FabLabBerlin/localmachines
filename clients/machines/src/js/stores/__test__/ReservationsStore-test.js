jest.dontMock('../../actionTypes');
jest.dontMock('../../getters');
jest.dontMock('lodash');
jest.dontMock('../MachineStore.js');
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
var MachineStore = require('../MachineStore');
var moment = require('moment');
var ReservationRulesStore = require('../ReservationRulesStore');
var ReservationsStore = require('../ReservationsStore');
var reactor = require('../../reactor');


function machines() {
  return {
    1: {
      Id: 1,
      Name: 'Printer123'
    }
  };
}

describe('ReservationsStore', function() {
  reactor.registerStores({
    machineStore: MachineStore,
    reservationRulesStore: ReservationRulesStore,
    reservationsStore: ReservationsStore
  });

  reactor.dispatch(actionTypes.SET_RESERVATION_RULES, []);

  describe('CREATE_SET_DATE', function() {
    it('sets the date and possible times at that day', function() {
      reactor.dispatch(actionTypes.SET_MACHINE_INFO, {
        machineInfo: machines()
      });
      reactor.dispatch(actionTypes.SET_RESERVATIONS, {
        reservations: []
      });
      reactor.dispatch(actionTypes.CREATE_EMPTY);
      reactor.dispatch(actionTypes.CREATE_SET_DATE, {
        date: moment('2015-10-16')
      });
      var possibleTimes = reactor.evaluateToJS(getters.getNewReservationTimes);
      expect(possibleTimes.length).toEqual(18);
      expect(_.first(possibleTimes).start.format('HH:mm')).toEqual('10:00');
      expect(_.first(possibleTimes).end.format('HH:mm')).toEqual('10:30');
      expect(_.last(possibleTimes).start.format('HH:mm')).toEqual('18:30');
      expect(_.last(possibleTimes).end.format('HH:mm')).toEqual('19:00');
      _.each(possibleTimes, (t) => {
        expect(t.availableMachineIds).toEqual([1]);
      });
    });
  });
});

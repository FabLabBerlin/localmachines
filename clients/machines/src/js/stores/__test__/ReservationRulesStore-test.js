jest.dontMock('lodash');
jest.dontMock('moment');
jest.dontMock('nuclear-js');
jest.dontMock('../../actionTypes');
jest.dontMock('../../getters');
jest.dontMock('../../reactor');
jest.dontMock('../MachineStore');
jest.dontMock('../ReservationRulesStore');
jest.dontMock('../ReservationsStore');
var actionTypes = require('../../actionTypes');
var getters = require('../../getters');
var moment = require('moment');
var reactor = require('../../reactor');


function existingReservations() {
  return [
    {  
      'Id': 1,
      'MachineId': 3,
      'UserId': 19,
      'TimeStart': '2015-10-15T11:30:00+02:00',
      'TimeEnd': '2015-10-15T18:00:00+02:00',
      'Created': '2015-10-02T17:29:51+02:00'
    },
    {  
      'Id': 2,
      'MachineId': 3,
      'UserId': 19,
      'TimeStart': '2015-10-06T14:00:00+02:00',
      'TimeEnd': '2015-10-06T15:00:00+02:00',
      'Created': '2015-10-05T11:23:58+02:00'
    }
  ];
}

function reservationRules() {
  return [
    {  
      'MachineId': 3,
      'Created': '2015-10-05T14:28:51+02:00',
      'Available': false,
      'Name': 'Lasercutter Introduction',
      'Tuesday': true,
      'DateEnd': '',
      'TimeEnd': '19:00',
      'Unavailable': true,
      'Monday': false,
      'Sunday': false,
      'DateStart': '',
      'TimeStart': '17:30',
      'Wednesday': false,
      'TimeZone': '',
      'Saturday': false,
      'Thursday': false,
      'Id': 1,
      'Friday': false
    },
    {  
      'MachineId': 2,
      'Created': '2015-10-05T14:47:24+02:00',
      'Available': false,
      'Name': '3D Printing Introduction',
      'Tuesday': true,
      'DateEnd': '',
      'TimeEnd': '19:00',
      'Unavailable': true,
      'Monday': false,
      'Sunday': false,
      'DateStart': '',
      'TimeStart': '17:30',
      'Wednesday': false,
      'TimeZone': '',
      'Saturday': false,
      'Thursday': false,
      'Id': 3,
      'Friday': false
    },
    {  
      'MachineId': 11,
      'Created': '2015-10-05T14:48:22+02:00',
      'Available': false,
      'Name': 'CNC-Milling Introduction',
      'Tuesday': false,
      'DateEnd': '2015-10-10',
      'TimeEnd': '18:00',
      'Unavailable': true,
      'Monday': false,
      'Sunday': false,
      'DateStart': '2015-10-10',
      'TimeStart': '09:00',
      'Wednesday': false,
      'TimeZone': '',
      'Saturday': false,
      'Thursday': false,
      'Id': 4,
      'Friday': false
    },
    {  
      'MachineId': 10,
      'Created': '2015-10-05T14:51:02+02:00',
      'Available': false,
      'Name': 'Christmas Presents',
      'Tuesday': true,
      'DateEnd': '2015-12-24',
      'TimeEnd': '23:59',
      'Unavailable': true,
      'Monday': true,
      'Sunday': true,
      'DateStart': '2015-12-20',
      'TimeStart': '17:00',
      'Wednesday': true,
      'TimeZone': '',
      'Saturday': true,
      'Thursday': true,
      'Id': 5,
      'Friday': true
    },
    {
      'MachineId': 10,
      'Created': '2015-10-05T14:51:02+02:00',
      'Available': true,
      'Name': 'Free i3 Printing Days',
      'Tuesday': true,
      'DateEnd': '2015-12-21',
      'TimeEnd': '',
      'Unavailable': false,
      'Monday': true,
      'Sunday': true,
      'DateStart': '2015-12-14',
      'TimeStart': '',
      'Wednesday': true,
      'TimeZone': '',
      'Saturday': true,
      'Thursday': true,
      'Id': 6,
      'Friday': true
    },
    {
      'MachineId': 3,
      'Created': '2015-10-05T14:51:02+02:00',
      'Available': false,
      'Name': 'Rule with no times/dates',
      'Tuesday': false,
      'DateEnd': '',
      'TimeEnd': '',
      'Unavailable': true,
      'Monday': false,
      'Sunday': false,
      'DateStart': '',
      'TimeStart': '',
      'Wednesday': false,
      'TimeZone': '',
      'Saturday': false,
      'Thursday': false,
      'Id': 7,
      'Friday': false
    }
  ];
}

function getMachineInfo() {
  return [
    {
      'Id': 3,
      'Name': 'Laser Cutter - Epilog Zing 6030 [€0.80/min]',
      'Shortname': 'ZLC',
      'Description': 'Cuts wood, plastic, paper. Fast.',
      'Image': '',
      'Available': true,
      'UnavailMsg': '',
      'UnavailTill': '0001-01-01T00:00:00Z',
      'Price': 0.8,
      'PriceUnit': 'minute',
      'Comments': 'asd',
      'Visible': true,
      'ConnectedMachines': '',
      'SwitchRefCount': 0,
      'UnderMaintenance': false
    },
    {
      'Id': 10,
      'Name': '3D Printer - 6 Honey Bunny (I3 Berlin) [€0.10/min]',
      'Shortname': 'I3B2',
      'Description': 'i3Berlin 3D Printer.',
      'Image': '',
      'Available': true,
      'UnavailMsg': '',
      'UnavailTill': '0001-01-01T00:00:00Z',
      'Price': 0.1,
      'PriceUnit': 'minute',
      'Comments': '',
      'Visible': true,
      'ConnectedMachines': '',
      'SwitchRefCount': 0,
      'UnderMaintenance': false
    },
    {
      'Id': 8,
      'Name': '3D Printer - 5 Pumpkin (I3 Berlin) [€0.10/min]',
      'Shortname': 'I3B1',
      'Description': 'i3Berlin 3D Printer',
      'Image': '',
      'Available': true,
      'UnavailMsg': '',
      'UnavailTill': '0001-01-01T00:00:00Z',
      'Price': 0.1,
      'PriceUnit': 'minute',
      'Comments': '',
      'Visible': true,
      'ConnectedMachines': '',
      'SwitchRefCount': 0,
      'UnderMaintenance': false
    }
  ];
}


describe('ReservationRulesStore', function() {
  var MachineStore = require('../MachineStore');
  var ReservationRulesStore = require('../ReservationRulesStore');
  var ReservationsStore = require('../ReservationsStore');

  reactor.registerStores({
    machineStore: MachineStore,
    reservationRulesStore: ReservationRulesStore,
    reservationsStore: ReservationsStore
  });

  var machineInfo = getMachineInfo();
  var reservations = existingReservations();
  reactor.dispatch(actionTypes.SET_MACHINE_INFO, { machineInfo });
  reactor.dispatch(actionTypes.SET_RESERVATIONS, { reservations });
  reactor.dispatch(actionTypes.SET_RESERVATION_RULES, reservationRules());

  it('works for the Lasercutting workshop every Tuesday', function() {
    reactor.dispatch(actionTypes.CREATE_EMPTY);
    var mid = 3;
    reactor.dispatch(actionTypes.CREATE_SET_MACHINE, { mid });
    var date = moment('2015-10-13');
    reactor.dispatch(actionTypes.CREATE_SET_DATE, { date });
    var tuesdayTimes = reactor.evaluateToJS(getters.getNewReservationTimes);

    var TIME_FMT = 'HH:mm';
    // Lasercutter not available on Tuesday 17:30 - 19:00 during workshop
    expect(tuesdayTimes.length).toEqual(18);
    expect(tuesdayTimes[0].start.format(TIME_FMT)).toEqual('10:00');
    expect(tuesdayTimes[0].end.format(TIME_FMT)).toEqual('10:30');
    expect(tuesdayTimes[0].availableMachineIds).toEqual([3, 8, 10]);

    expect(tuesdayTimes[14].start.format(TIME_FMT)).toEqual('17:00');
    expect(tuesdayTimes[14].end.format(TIME_FMT)).toEqual('17:30');
    expect(tuesdayTimes[14].availableMachineIds).toEqual([3, 8, 10]);

    expect(tuesdayTimes[15].start.format(TIME_FMT)).toEqual('17:30');
    expect(tuesdayTimes[15].end.format(TIME_FMT)).toEqual('18:00');
    expect(tuesdayTimes[15].availableMachineIds).toEqual([8, 10]);

    expect(tuesdayTimes[16].start.format(TIME_FMT)).toEqual('18:00');
    expect(tuesdayTimes[16].end.format(TIME_FMT)).toEqual('18:30');
    expect(tuesdayTimes[16].availableMachineIds).toEqual([8, 10]);

    expect(tuesdayTimes[17].start.format(TIME_FMT)).toEqual('18:30');
    expect(tuesdayTimes[17].end.format(TIME_FMT)).toEqual('19:00');
    expect(tuesdayTimes[17].availableMachineIds).toEqual([8, 10]);

    // Lasercutter available whole Wednesday
    date = moment('2015-10-14');
    reactor.dispatch(actionTypes.CREATE_SET_DATE, { date });
    var wednesdayTimes = reactor.evaluateToJS(getters.getNewReservationTimes);
    expect(wednesdayTimes.length).toEqual(18);

    for (var i = 0; i < 18; i++) {
      expect(wednesdayTimes[i].availableMachineIds).toEqual([3, 8, 10]);
    }
  });

  it('works for availability overlapping an unavailability', function() {
    reactor.dispatch(actionTypes.CREATE_EMPTY);
    var mid = 3;
    reactor.dispatch(actionTypes.CREATE_SET_MACHINE, { mid });
    var i;

    // i3 not available December Wednesday, 23rd
    var date = moment('2015-12-23');
    reactor.dispatch(actionTypes.CREATE_SET_DATE, { date });
    var times = reactor.evaluateToJS(getters.getNewReservationTimes);
    expect(times.length).toEqual(18);
    for (i = 0; i < 18; i++) {
      expect(times[i].availableMachineIds).toEqual([3, 8]);
    }

    // i3 not available December Tuesday, 22nd
    date = moment('2015-12-22');
    reactor.dispatch(actionTypes.CREATE_SET_DATE, { date });
    times = reactor.evaluateToJS(getters.getNewReservationTimes);
    expect(times.length).toEqual(18);
    for (i = 0; i < 18; i++) {
      if (i < 15) {
        expect(times[i].availableMachineIds).toEqual([3, 8]);
      } else {
        // Lasercutter workshop in the evening...
        expect(times[i].availableMachineIds).toEqual([8]);
      }
    }

    // i3 available December Saturday, 19th because of Free i3 Printing days
    date = moment('2015-12-19');
    reactor.dispatch(actionTypes.CREATE_SET_DATE, { date });
    times = reactor.evaluateToJS(getters.getNewReservationTimes);
    expect(times.length).toEqual(12);
    for (i = 0; i < 12; i++) {
      expect(times[i].availableMachineIds).toEqual([3, 8, 10]);
    }
  });

  it('ignores rules with neither time/date specified', function() {
    reactor.dispatch(actionTypes.CREATE_EMPTY);
    var mid = 3;
    reactor.dispatch(actionTypes.CREATE_SET_MACHINE, { mid });
    var i;

    // Everything is available on November Wednesday, 11th
    var date = moment('2015-11-11');
    reactor.dispatch(actionTypes.CREATE_SET_DATE, { date });
    var times = reactor.evaluateToJS(getters.getNewReservationTimes);
    expect(times.length).toEqual(18);
    for (i = 0; i < 18; i++) {
      expect(times[i].availableMachineIds).toEqual([3, 8, 10]);
    }
  });

  it('takes existing reservations into account', function() {
    reactor.dispatch(actionTypes.CREATE_EMPTY);
    var mid = 3;
    reactor.dispatch(actionTypes.CREATE_SET_MACHINE, { mid });

    var date = moment('2015-10-15');
    reactor.dispatch(actionTypes.CREATE_SET_DATE, { date });
    var times = reactor.evaluateToJS(getters.getNewReservationTimes);
    expect(times.length).toEqual(18);
    for (var i = 0; i < 18; i++) {
      if (i < 3 || i > 15) {
        expect(times[i].availableMachineIds).toEqual([3, 8, 10]);
      } else {
        expect(times[i].availableMachineIds).toEqual([8, 10]);
      }
    }
  });
});

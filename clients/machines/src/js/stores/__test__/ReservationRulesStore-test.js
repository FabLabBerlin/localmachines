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
      'MachineId': 12,
      'Created': '2015-10-05T14:51:02+02:00',
      'Available': false,
      'Name': 'Christmas Presents',
      'Tuesday': false,
      'DateEnd': '2015-12-24',
      'TimeEnd': '23:59',
      'Unavailable': true,
      'Monday': false,
      'Sunday': false,
      'DateStart': '2015-12-20',
      'TimeStart': '17:00',
      'Wednesday': false,
      'TimeZone': '',
      'Saturday': false,
      'Thursday': false,
      'Id': 5,
      'Friday': false
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

  it('makes getNewReservationTimes consider the rules', function() {
    var machineInfo = [
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
    reactor.dispatch(actionTypes.SET_MACHINE_INFO, { machineInfo });
    reactor.dispatch(actionTypes.SET_RESERVATION_RULES, reservationRules());
    reactor.dispatch(actionTypes.CREATE_EMPTY);
    var mid = 3;
    reactor.dispatch(actionTypes.CREATE_SET_MACHINE, { mid });
    var date = moment('2015-10-13');
    reactor.dispatch(actionTypes.CREATE_SET_DATE, { date });
    var tuesdayTimes = reactor.evaluateToJS(getters.getNewReservationTimes);

    var TIME_FMT = 'HH:mm';
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
  });
});

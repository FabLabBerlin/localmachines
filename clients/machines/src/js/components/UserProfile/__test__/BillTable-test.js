jest.dontMock('lodash');
jest.dontMock('moment');
jest.dontMock('nuclear-js');
jest.dontMock('react');
jest.dontMock('../../../actionTypes');
jest.dontMock('../../../getters');
jest.dontMock('../../../reactor');
jest.dontMock('../BillTable');
jest.dontMock('../helpers');
jest.dontMock('../../../stores/UserStore.js');


var React = require('react');
var actionTypes = require('../../../actionTypes');
var getters = require('../../../getters');
var BillTable = React.createFactory(require('../BillTable'));
var Nuclear = require('nuclear-js');
var reactor = require('../../../reactor');
var toImmutable = Nuclear.toImmutable;


var state = {
  'userStore': {
    'userId': 0,
    'isLogged': false,
    'firstTry': true,
    'bill': {
      'User': {
        'Id': 19,
        'FirstName': 'Philip',
        'LastName': 'Silva',
        'Username': 'psilva',
        'Email': 'psilva.de@gmail.com',
        'InvoiceAddr': 'There is one',
        'ShipAddr': '0',
        'ClientId': 696,
        'B2b': false,
        'Company': '',
        'VatUserId': '',
        'VatRate': 0,
        'UserRole': 'admin',
        'Created': '2015-06-04T08:34:51+02:00',
        'Comments': 'Fastbill Kd-Nr. 696',
        'Phone': '0002',
        'ZipCode': '',
        'City': '',
        'CountryCode': ''
      },
      'Purchases': {
        'Data': [
          {
            'Id': 8353,
            'Type': 'activation',
            'ProductId': 0,
            'Created': '0001-01-01T00:00:00Z',
            'UserId': 19,
            'TimeStart': '2015-12-26T10:42:49+01:00',
            'TimeEnd': '2015-12-26T10:44:01+01:00',
            'Quantity': 1.2,
            'PricePerUnit': 0.8,
            'PriceUnit': 'minute',
            'Vat': 0,
            'Cancelled': false,
            'TotalPrice': 0.96,
            'DiscountedTotal': 0.96,
            'ActivationRunning': false,
            'ReservationDisabled': false,
            'Machine': {
              'Id': 3,
              'Name': 'Laser Cutter - Epilog Zing 6030',
              'Shortname': 'ZLC',
              'Description': 'Cuts wood, plastic, paper. Fast.',
              'Image': 'machine-3.svg',
              'Available': true,
              'UnavailMsg': '',
              'UnavailTill': '0001-01-01T00:00:00Z',
              'Price': 0.8,
              'PriceUnit': 'minute',
              'Comments': 'asd',
              'Visible': true,
              'ConnectedMachines': '',
              'SwitchRefCount': 3,
              'UnderMaintenance': false,
              'ReservationPriceStart': null,
              'ReservationPriceHourly': 5
            },
            'MachineId': 3,
            'Activation': null,
            'Reservation': null,
            'MachineUsage': 0,
            'Memberships': [
              {
                'Id': 4,
                'Title': '3D Print Club',
                'ShortName': '3DPC',
                'DurationMonths': 12,
                'MonthlyPrice': 10,
                'MachinePriceDeduction': 100,
                'AffectedMachines': '[2,4,10,7,6,12,8,13,31,29,34]',
                'AutoExtend': true,
                'AutoExtendDurationMonths': 1
              }
            ]
          },
          {
            'Id': 8354,
            'Type': 'activation',
            'ProductId': 0,
            'Created': '0001-01-01T00:00:00Z',
            'UserId': 19,
            'TimeStart': '2015-12-26T10:46:35+01:00',
            'TimeEnd': '2015-12-26T10:47:27+01:00',
            'Quantity': 0.866666666666666,
            'PricePerUnit': 0.1,
            'PriceUnit': 'minute',
            'Vat': 0,
            'Cancelled': false,
            'TotalPrice': 0.08666666666666661,
            'DiscountedTotal': 0,
            'ActivationRunning': false,
            'ReservationDisabled': false,
            'Machine': {
              'Id': 8,
              'Name': '3D Printer - 5 Pumpkin (I3 Berlin)',
              'Shortname': 'I3B1',
              'Description': 'i3Berlin 3D Printer',
              'Image': 'machine-8.svg',
              'Available': true,
              'UnavailMsg': '',
              'UnavailTill': '0001-01-01T00:00:00Z',
              'Price': 0.1,
              'PriceUnit': 'minute',
              'Comments': '',
              'Visible': true,
              'ConnectedMachines': '',
              'SwitchRefCount': 0,
              'UnderMaintenance': false,
              'ReservationPriceStart': null,
              'ReservationPriceHourly': null
            },
            'MachineId': 8,
            'Activation': null,
            'Reservation': null,
            'MachineUsage': 0,
            'Memberships': [
              {
                'Id': 4,
                'Title': '3D Print Club',
                'ShortName': '3DPC',
                'DurationMonths': 12,
                'MonthlyPrice': 10,
                'MachinePriceDeduction': 100,
                'AffectedMachines': '[2,4,10,7,6,12,8,13,31,29,34]',
                'AutoExtend': true,
                'AutoExtendDurationMonths': 1
              }
            ]
          },
          {
            'Id': 8362,
            'Type': 'reservation',
            'ProductId': 0,
            'Created': '2015-12-26T11:16:48+01:00',
            'UserId': 19,
            'TimeStart': '2015-12-27T13:00:00+01:00',
            'TimeEnd': '2015-12-27T13:30:00+01:00',
            'Quantity': 1,
            'PricePerUnit': 2.5,
            'PriceUnit': '30 minutes',
            'Vat': 0,
            'Cancelled': false,
            'TotalPrice': 2.5,
            'DiscountedTotal': 2.5,
            'ActivationRunning': false,
            'ReservationDisabled': false,
            'Machine': {
              'Id': 11,
              'Name': 'CNC Router',
              'Shortname': 'CNC',
              'Description': '',
              'Image': 'machine-11.svg',
              'Available': true,
              'UnavailMsg': '',
              'UnavailTill': '0001-01-01T00:00:00Z',
              'Price': 0.8,
              'PriceUnit': 'minute',
              'Comments': '',
              'Visible': true,
              'ConnectedMachines': '',
              'SwitchRefCount': 0,
              'UnderMaintenance': false,
              'ReservationPriceStart': null,
              'ReservationPriceHourly': 5
            },
            'MachineId': 11,
            'Activation': null,
            'Reservation': null,
            'MachineUsage': 0,
            'Memberships': [
              {
                'Id': 4,
                'Title': '3D Print Club',
                'ShortName': '3DPC',
                'DurationMonths': 12,
                'MonthlyPrice': 10,
                'MachinePriceDeduction': 100,
                'AffectedMachines': '[2,4,10,7,6,12,8,13,31,29,34]',
                'AutoExtend': true,
                'AutoExtendDurationMonths': 1
              }
            ]
          }
        ]
      }
    },
    'memberships': [
      {
        'Id': 62,
        'UserId': 261,
        'MembershipId': 4,
        'StartDate': '2015-08-05T02:00:00+02:00',
        'EndDate': '2016-08-05T02:00:00+02:00',
        'Title': '3D Print Club',
        'ShortName': '3DPC',
        'DurationMonths': 12,
        'Unit': '',
        'MonthlyPrice': 10,
        'MachinePriceDeduction': 100,
        'AffectedMachines': '[2,4,6,10,8,7,12,13]',
        'AutoExtend': true
      }
    ]
  }
};


describe('BillTable', function() {
  var $ = require('jquery');
  var UserStore = require('../../../stores/UserStore');

  reactor.registerStores({
    userStore: UserStore
  });

  describe('render', function() {
    it('renders the activations and the totals', function() {
      var data = state.userStore.bill;
      reactor.dispatch(actionTypes.SET_BILL, { data });
      data = state.userStore.memberships;
      reactor.dispatch(actionTypes.SET_MEMBERSHIPS, { data });

      var billTable = new BillTable();
      var html = React.renderToString(billTable);
      /* Activations */
      expect(html).toContain('Laser Cutter - Epilog Zing 6030');
      expect(html).toContain('1m 12s');
      expect(html).toContain('0.96');
      expect(html).toContain('3D Printer - 5 Pumpkin (I3 Berlin)');
      expect(html).toContain('51s');
      expect(html).toContain('0.00');
      /* Reservations */
      expect(html).toContain('CNC Router (Reservation)');
      expect(html).toContain('30m 0s');
      expect(html).toContain('2.50');
      /* Totals */
      expect(html).toContain('Total Pay-As-You-Go');
      expect(html).toContain('3.46');
      expect(html).toContain('Total Memberships');
      expect(html).toContain('10.00');
      expect(html).toContain('Total');
      expect(html).toContain('13.46');
    });
  });
});

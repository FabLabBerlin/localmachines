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
    'billInfo': {
      'User': {
        'Id': 261,
        'FirstName': 'bla',
        'LastName': 'bla',
        'Username': 'bla',
        'Email': 'bla@example.com',
        'InvoiceAddr': '',
        'ShipAddr': '',
        'ClientId': 0,
        'B2b': false,
        'Company': '',
        'VatUserId': '',
        'VatRate': 0,
        'UserRole': '',
        'Created': '2015-10-07T10:30:21+02:00',
        'Comments': '',
        'Phone': '',
        'ZipCode': '',
        'City': '',
        'CountryCode': ''
      },
      'Purchases': [
        {
          'Activation': {
            'Id': 3195,
            'InvoiceId': 0,
            'UserId': 261,
            'MachineId': 8,
            'Active': false,
            'TimeStart': '2015-10-07T10:32:59+02:00',
            'TimeEnd': '2015-10-07T10:33:51+02:00',
            'TimeTotal': 51,
            'UsedKwh': 0,
            'DiscountPercents': 0,
            'DiscountFixed': 0,
            'VatRate': 0,
            'CommentRef': '',
            'Invoiced': false,
            'Changed': false
          },
          'Machine': {
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
          },
          'MachineUsage': 0.85,
          'User': {
            'Id': 261,
            'FirstName': 'bla',
            'LastName': 'bla',
            'Username': 'bla',
            'Email': 'bla@example.com',
            'InvoiceAddr': '',
            'ShipAddr': '',
            'ClientId': 0,
            'B2b': false,
            'Company': '',
            'VatUserId': '',
            'VatRate': 0,
            'UserRole': '',
            'Created': '2015-10-07T10:30:21+02:00',
            'Comments': '',
            'Phone': '',
            'ZipCode': '',
            'City': '',
            'CountryCode': ''
          },
          'Memberships': [
            {
              'Id': 4,
              'Title': '3D Print Club',
              'ShortName': '3DPC',
              'DurationMonths': 12,
              'MonthlyPrice': 10,
              'MachinePriceDeduction': 100,
              'AffectedMachines': '[2,4,6,10,8,7,12,13]',
              'AutoExtend': true,
              'AutoExtendDurationMonths': 1
            }
          ],
          'TotalPrice': 0.085,
          'DiscountedTotal': 0
        },
        {
          'Activation': {
            'Id': 3196,
            'InvoiceId': 0,
            'UserId': 261,
            'MachineId': 3,
            'Active': false,
            'TimeStart': '2015-10-07T10:33:54+02:00',
            'TimeEnd': '2015-10-07T10:35:06+02:00',
            'TimeTotal': 72,
            'UsedKwh': 0,
            'DiscountPercents': 0,
            'DiscountFixed': 0,
            'VatRate': 0,
            'CommentRef': '',
            'Invoiced': false,
            'Changed': false
          },
          'Machine': {
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
          'MachineUsage': 1.2,
          'User': {
            'Id': 261,
            'FirstName': 'bla',
            'LastName': 'bla',
            'Username': 'bla',
            'Email': 'bla@example.com',
            'InvoiceAddr': '',
            'ShipAddr': '',
            'ClientId': 0,
            'B2b': false,
            'Company': '',
            'VatUserId': '',
            'VatRate': 0,
            'UserRole': '',
            'Created': '2015-10-07T10:30:21+02:00',
            'Comments': '',
            'Phone': '',
            'ZipCode': '',
            'City': '',
            'CountryCode': ''
          },
          'Memberships': [
            {
              'Id': 4,
              'Title': '3D Print Club',
              'ShortName': '3DPC',
              'DurationMonths': 12,
              'MonthlyPrice': 10,
              'MachinePriceDeduction': 100,
              'AffectedMachines': '[2,4,6,10,8,7,12,13]',
              'AutoExtend': true,
              'AutoExtendDurationMonths': 1
            }
          ],
          'TotalPrice': 0.96,
          'DiscountedTotal': 0.96
        }
      ]
    },
    'membershipInfo': [
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
      var data = state.userStore.billInfo;
      reactor.dispatch(actionTypes.SET_BILL_INFO, { data });
      data = state.userStore.membershipInfo;
      reactor.dispatch(actionTypes.SET_MEMBERSHIP_INFO, { data });

      var billTable = new BillTable();
      var html = React.renderToString(billTable);
      /* Activations */
      expect(html).toContain('Laser Cutter - Epilog Zing 6030 [€0.80/min]');
      expect(html).toContain('1m 12s');
      expect(html).toContain('0.96');
      expect(html).toContain('3D Printer - 5 Pumpkin (I3 Berlin) [€0.10/min]');
      expect(html).toContain('52s');
      expect(html).toContain('0.00');
      /* Totals */
      expect(html).toContain('Total Pay-As-You-Go');
      expect(html).toContain('0.96');
      expect(html).toContain('Total Memberships');
      expect(html).toContain('10.00');
      expect(html).toContain('Total');
      expect(html).toContain('10.96');
    });
  });
});

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
      'UserId': 1,
      'UserClientId': 0,
      'UserFirstName': 'Regular',
      'UserLastName': 'User',
      'Username': 'user',
      'UserEmail': 'user@example.com',
      'DebitorNumber': 'Undefined',
      'Activations': [
        {
          'MachineId': 1,
          'MachineName': 'Laydrop 3D Printer',
          'MachineProductId': 'Undefined',
          'MachineUsage': 0.01,
          'MachineUsageUnit': 'hour',
          'MachinePricePerUnit': 16,
          'UserId': 1,
          'UserClientId': 0,
          'UserFirstName': 'Regular',
          'UserLastName': 'User',
          'Username': 'user',
          'UserEmail': 'user@example.com',
          'UserDebitorNumber': 'Undefined',
          'UserInvoiceAddr': '0',
          'UserZipCode': '',
          'UserCity': '',
          'UserCountryCode': '',
          'UserPhone': '',
          'UserComments': '',
          'Memberships': [
            {
              'Id': 0,
              'Title': '1 Month Basic',
              'ShortName': 'SMB',
              'Duration': 30,
              'Unit': 'days',
              'MonthlyPrice': 0,
              'MachinePriceDeduction': 50,
              'AffectedMachines': '[1,2,3]'
            }
          ],
          'TotalPrice': 0.16,
          'DiscountedTotal': 0.08,
          'TimeStart': '2015-09-11T10:01:44+02:00',
          'TimeEnd': '2015-09-11T10:01:56+02:00'
        },
        {
          'MachineId': 2,
          'MachineName': 'MakerBot 3D Printer',
          'MachineProductId': 'Undefined',
          'MachineUsage': 0.01,
          'MachineUsageUnit': 'hour',
          'MachinePricePerUnit': 16,
          'UserId': 1,
          'UserClientId': 0,
          'UserFirstName': 'Regular',
          'UserLastName': 'User',
          'Username': 'user',
          'UserEmail': 'user@example.com',
          'UserDebitorNumber': 'Undefined',
          'UserInvoiceAddr': '0',
          'UserZipCode': '',
          'UserCity': '',
          'UserCountryCode': '',
          'UserPhone': '',
          'UserComments': '',
          'Memberships': [
            {
              'Id': 0,
              'Title': '1 Month Basic',
              'ShortName': 'SMB',
              'Duration': 30,
              'Unit': 'days',
              'MonthlyPrice': 0,
              'MachinePriceDeduction': 50,
              'AffectedMachines': '[1,2,3]'
            }
          ],
          'TotalPrice': 0.16,
          'DiscountedTotal': 0.08,
          'TimeStart': '2015-09-11T10:02:05+02:00',
          'TimeEnd': '2015-09-11T10:02:11+02:00'
        }
      ],
      'InvoiceAddr': '0',
      'ZipCode': '',
      'City': '',
      'CountryCode': '',
      'Phone': '',
      'Comments': ''
    },
    'membershipInfo': [
      {
        'Id': 1,
        'UserId': 1,
        'MembershipId': 2,
        'StartDate': '2015-09-01T00:00:00+02:00',
        'Title': '1 Month Basic',
        'ShortName': 'SMB',
        'Duration': 30,
        'Unit': 'days',
        'MonthlyPrice': 100,
        'MachinePriceDeduction': 50,
        'AffectedMachines': '[1,2,3]'
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
      expect(html).toContain('MakerBot 3D Printer');
      expect(html).toContain('6 s');
      expect(html).toContain('0.08');
      expect(html).toContain('Laydrop 3D Printer');
      expect(html).toContain('12 s');
      expect(html).toContain('0.08');
      /* Totals */
      expect(html).toContain('Total Pay-As-You-Go');
      expect(html).toContain('0.16');
      expect(html).toContain('Total Memberships');
      expect(html).toContain('100.00');
      expect(html).toContain('Total');
      expect(html).toContain('100.16');
    });
  });
});

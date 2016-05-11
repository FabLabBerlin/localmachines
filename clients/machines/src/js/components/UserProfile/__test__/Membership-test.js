jest.dontMock('lodash');
jest.dontMock('moment');
jest.dontMock('nuclear-js');
jest.dontMock('react');
jest.dontMock('../../../actionTypes');
jest.dontMock('../../../getters');
jest.dontMock('../helpers');
jest.dontMock('../Membership');


var React = require('react');
var actionTypes = require('../../../actionTypes');
var getters = require('../../../getters');
var Membership = React.createFactory(require('../Membership'));
var Nuclear = require('nuclear-js');
var reactor = require('../../../reactor');
var toImmutable = Nuclear.toImmutable;


var state = {
  'userStore': {
    'memberships': toImmutable([
      {
        'Id': 1,
        'UserId': 1,
        'MembershipId': 2,
        'StartDate': '2015-09-01T00:00:00+02:00',
        'EndDate': '2015-10-01T00:00:00+02:00',
        'Title': '1 Month Basic',
        'ShortName': 'SMB',
        'DurationMonths': 1,
        'Unit': 'days',
        'MonthlyPrice': 100,
        'MachinePriceDeduction': 50,
        'AffectedMachines': '[1,2,3]'
      }
    ])
  }
};


describe('Membership', function() {
  describe('render', function() {
    it('renders the memberships and the totals', function() {
      var membership = new Membership({
        memberships: state.userStore.memberships
      });
      var html = React.renderToString(membership);
      /* Memberships */
      expect(html).toContain('Name');
      expect(html).toContain('1 Month Basic');
      expect(html).toContain('Start Date');
      expect(html).toContain('01. Sep 2015');
      expect(html).toContain('End Date');
      expect(html).toContain('01. Oct 2015');
      /* Totals */
      expect(html).toContain('Total/month');
      expect(html).toContain('100.00');
    });
  });
});

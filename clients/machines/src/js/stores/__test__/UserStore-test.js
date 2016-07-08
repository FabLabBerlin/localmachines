jest.dontMock('../../actionTypes');
jest.dontMock('../../getters');
jest.dontMock('../../modules/Machines/actionTypes');
jest.dontMock('../../modules/Machines/index');
jest.dontMock('../../modules/Machines');
jest.dontMock('../../modules/Machines/stores/store');
jest.dontMock('nuclear-js');
jest.dontMock('../../reactor');
jest.dontMock('../UserStore');

var actionTypes = require('../../actionTypes');
var getters = require('../../getters');
var reactor = require('../../reactor');


function getUser() {
  return {
    FirstName: 'Regular',
    Id: 2,
    LastName: 'Admin',
    UserRole: 'admin'
  };
}

describe('UserStore', function() {
  var UserStore = require('../UserStore');

  reactor.registerStores({
    userStore: UserStore
  });

  describe('SET_USER', function() {
    it('does changes visible via getUser', function() {
      var user = getUser();
      reactor.dispatch(actionTypes.SET_USER, { user });
      var actual = reactor.evaluateToJS(getters.getUser);
      expect(actual).toEqual(getUser());
    });
  });
});

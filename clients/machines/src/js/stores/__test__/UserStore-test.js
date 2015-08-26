jest.dontMock('../../actionTypes');
jest.dontMock('../../getters');
jest.dontMock('nuclear-js');
jest.dontMock('../../reactor');
jest.dontMock('../UserStore');

var actionTypes = require('../../actionTypes');
var getters = require('../../getters');
var reactor = require('../../reactor');


function user() {
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

  describe('SET_USER_INFO', function() {
    it('does changes visible via getUserInfo', function() {
      var userInfo = user();
      reactor.dispatch(actionTypes.SET_USER_INFO, { userInfo });
      var actual = reactor.evaluateToJS(getters.getUserInfo);
      expect(actual).toEqual(user());
    });
  });
});

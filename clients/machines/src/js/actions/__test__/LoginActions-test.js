jest.dontMock('../LoginActions');
jest.mock('jquery');

var $ = require('jquery');
var LoginActions = require('../LoginActions');

describe('LoginActions', function() {
  describe('submitLoginForm', function() {
    it('POSTs to /api/users/login', function() {
      LoginActions.submitLoginForm({
        username: 'foo',
        password: 'bar'
      });
      expect($.ajax).toBeCalledWith({
        url: '/api/users/login',
        dataType: 'json',
        type: 'POST',
        data: {
          username: 'foo',
          password: 'bar'
        },
        success: jasmine.any(Function),
        error: jasmine.any(Function)
      });
    });
  });

  describe('logout', function() {
    it('GETs /api/users/logout', function() {
      LoginActions.logout();
      expect($.ajax).toBeCalledWith({
        url: '/api/users/logout',
        type: 'GET',
        cache: false,
        success: jasmine.any(Function),
        error: jasmine.any(Function)
      });
    });
  });
});

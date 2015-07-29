jest.mock('jquery');
jest.dontMock('../LoginStore.js');

describe('LoginStore', function() {
  var $ = require('jquery');
  var LoginStore = require('../LoginStore');


  describe('apiGetLogout', function() {
    it('sends a logout API request', function() {
      LoginStore.apiGetLogout();
      expect($.ajax).toBeCalledWith({
        cache: false,
        type: 'GET',
        url: '/api/users/logout',
        success: jasmine.any(Function),
        error: jasmine.any(Function)
      });
    });
  });

  describe('apiPostLogin', function() {
    it('sends a login POST request', function() {
      var loginInfo = {
        foo: 'bar'
      };
      LoginStore.apiPostLogin(loginInfo);
      expect($.ajax).toBeCalledWith({
        url: '/api/users/login',
        dataType: 'json',
        type: 'POST',
        data: loginInfo,
        success: jasmine.any(Function),
        error: jasmine.any(Function)
      });
    });
  });

  describe('apiPostLoginNFC', function() {
    it('sends a login(uid) POST request', function() {
      LoginStore.apiPostLoginNFC(123);
      expect($.ajax).toBeCalledWith({
        url: '/api/users/loginuid',
        method: 'POST',
        data: {
          uid: 123
        },
        success: jasmine.any(Function),
        error: jasmine.any(Function)
      });
    });
  });
});

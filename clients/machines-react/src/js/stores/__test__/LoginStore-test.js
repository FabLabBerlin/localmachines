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
});

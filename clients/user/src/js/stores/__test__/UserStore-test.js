jest.dontMock('../UserStore.js');

describe('UserStore test', function() {
  it('call into $.ajax with the correct params', function() {
    var $ = require('jquery');

    var UserStore = require('../UserStore');

    UserStore.logoutFromServer();

    expect($.ajax).toBeCalledWith({
      cache: false,
      dataType: 'json',
      type: 'GET',
      url: '/api/users/logout',
      success: jasmine.any(Function),
      error: jasmine.any(Function)
    });
  });
});

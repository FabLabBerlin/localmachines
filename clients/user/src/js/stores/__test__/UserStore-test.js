jest.dontMock('../UserStore.js');

describe('UserStore test', function() {

  /*
   * Test API calls
   * without params or only uid
   */
  it('call into $.ajax with no param or uid', function() {
    var $ = require('jquery');
    var UserStore = require('../UserStore');
    var uidTest = 3;

    /*
     * Test logout
     */
    UserStore.logoutFromServer();
    expect($.ajax).toBeCalledWith({
      cache: false,
      dataType: 'json',
      type: 'GET',
      url: '/api/users/logout',
      success: jasmine.any(Function),
      error: jasmine.any(Function)
    });

    /*
     * Test getUserStateFromServer
     */
    UserStore.getUserStateFromServer(uidTest);
    expect($.ajax).toBeCalledWith({
      url: '/api/users/3',
      dataType: 'json',
      type: 'GET',
      success: jasmine.any(Function),
      error: jasmine.any(Function)
    });

    /*
     * Test getMachineFromServer
     */
    UserStore.getMachineFromServer(uidTest);
    expect($.ajax).toBeCalledWith({
      url: '/api/users/3/machinepermissions',
      dataType: 'json',
      type: 'GET',
      success: jasmine.any(Function),
      error: jasmine.any(Function)
    });

    /*
     * Test getMembershipFromServer
     */
    UserStore.getMembershipFromServer(uidTest);
    expect($.ajax).toBeCalledWith({
      url: '/api/users/3/machinepermissions',
      dataType: 'json',
      type: 'GET',
      success: jasmine.any(Function),
      error: jasmine.any(Function)
    });
  });

});

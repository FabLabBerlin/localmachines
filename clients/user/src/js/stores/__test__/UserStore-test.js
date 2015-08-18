jest.dontMock('../UserStore.js');

describe('UserStore test', function() {
  /*
   * Dependencies
   */
  var $ = require('jquery');
  var UserStore = require('../UserStore');

  /*
   * Data used for testing
   */
  UserStore._state = {
    userId: 5,
    isLogged: true,
    rawInfoUser: {
      Id: 5,
      FirstName: "Regular",
      LastName: "User",
      Username: "user",
      Email: "user@example.com",
      InvoiceAddr: 0,
      ShipAddr: 0,
      ClientId: 0,
      B2b: false,
      Company: "",
      VatUserId: "",
      VatRate: 0,
      UserRole: "",
      Created: "0001-01-01T00:00:00Z",
      Comments: ""
    },
    rawInfoMachine: [
      {
        Id: 1,
        Name: "Laydrop 3D Printer",
        Shortname: "MB3DP",
        Description: "NYC 3D printer 4 real and 4 life.",
        Image: "",
        Available: true,
        UnavailMsg: "",
        UnavailTill: "0001-01-01T00:00:00Z",
        Price: 16,
        PriceUnit: "hour",
        Comments: "",
        Visible: true,
        ConnectedMachines: "[3]",
        SwitchRefCount: 0
      },
      {
        Id: 2,
        Name: "MakerBot 3D Printer",
        Shortname: "MB3DP",
        Description: "NYC 3D printer 4 real and 4 life.",
        Image: "machine-2.svg",
        Available: true,
        UnavailMsg: "",
        UnavailTill: "0001-01-01T00:00:00Z",
        Price: 16,
        PriceUnit: "hour",
        Comments: "",
        Visible: true,
        ConnectedMachines: "",
        SwitchRefCount: 0
      }
    ],
    rawInfoBill: {
      TotalTime: 257,
      TotalPrice: 1.1422222,
      Details: [
        {
          MachineId: 1,
          MachineName: "Laydrop 3D Printer",
          Price: "1.1422222",
          Time: 257
        },
      ]
    },
    rawInfoMembership: []
  };

  var testInfoMachine = UserStore._state['rawInfoMachine'];
  var uidTest = 5;

  var userUpdateFullState = {
    User: {
      Id: 5,
      FirstName: "user",
      LastName: "Update",
      Username: "State",
      Email: "user@example.com",
      InvoiceAddr: 0,
      ShipAddr: 0,
      ClientId: 0,
      B2b: false,
      Company: "",
      VatUserId: "",
      VatRate: 0,
      UserRole: "",
      Created: "0001-01-01T00:00:00Z",
      Comments: ""
    },
  };

  var userUpdateState = {
    FirstName: "user",
    LastName: "Update",
    Username: "State",
    Email: "user@example.com",
    InvoiceAddr: 0,
    ShipAddr: 0
  };

  var wantedResponse = {
    FirstName: "Regular",
    LastName: "User",
    Username: "user",
    Email: "user@example.com"
  };

  var emptyState = {
    userId : 0,
    isLogged : false,
    rawInfoUser : {},
    rawInfoMachine : [],
    rawInfoBill: {},
    rawInfoMembership : {}
  };

  var loginInfoTest = {
    username: 'test',
    password: 'test'
  };

  /*
   *
   * TEST API CALLS
   * without params or only uid
   *
   */

  /*
   * Test logout
   */
  it('test logoutFromServer', function() {
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

  it('test getUserStateFromServer', function(){
    /*
     * Test getUserStateFromServer
     */
    UserStore.getUserStateFromServer(uidTest);
    expect($.ajax).toBeCalledWith({
      url: '/api/users/5',
      dataType: 'json',
      type: 'GET',
      success: jasmine.any(Function),
      error: jasmine.any(Function)
    });
  })

  /*
   * Test getMachineFromServer
   */
  it('test getMachineFromServer', function() {
    UserStore.getMachineFromServer(uidTest);
    expect($.ajax).toBeCalledWith({
      url: '/api/users/5/machinepermissions',
      dataType: 'json',
      type: 'GET',
      success: jasmine.any(Function),
      error: jasmine.any(Function)
    });
  });

  /*
   * Test getInfoBillFromServer
   */
  it('test getInfoBillFromServer', function() {
    UserStore.getInfoBillFromServer(uidTest);
    expect($.ajax).toBeCalledWith({
      url: '/api/users/5/bill',
      dataType: 'json',
      type: 'GET',
      success: jasmine.any(Function),
      error: jasmine.any(Function)
    });
  });

  /*
   * Test getMembershipFromServer
   */
  it('test getMembershipFromServer', function() {
    UserStore.getMembershipFromServer(uidTest);
    expect($.ajax).toBeCalledWith({
      url: '/api/users/5/machinepermissions',
      dataType: 'json',
      type: 'GET',
      success: jasmine.any(Function),
      error: jasmine.any(Function)
    });
  });

  /*
   * Test submitLoginFormToServer
   */
  it('test submitLoginFormToServer', function() {
    UserStore.submitLoginFormToServer(loginInfoTest);
    expect($.ajax).toBeCalledWith({
      url: '/api/users/login',
      dataType: 'json',
      type: 'POST',
      data: loginInfoTest,
      success: jasmine.any(Function),
      error: jasmine.any(Function)
    });
  });

  /*
   * Test updatePassword
   */
  it('test updatePassword', function() {
    UserStore.updatePassword('passwordTest');
    expect($.ajax).toBeCalledWith({
      url: '/api/users/5/password',
      dataType: 'json',
      type: 'POST',
      data: {
        password: 'passwordTest'
      },
      success: jasmine.any(Function),
      error: jasmine.any(Function)
    });
  });

  /*
   * Test formatUserStateToSendToUserPage
   */
  it('test formatUserStateToSendToUserPage', function() {
    expect( UserStore.formatUserStateToSendToUserPage() ).toEqual(wantedResponse);
  });

  /*
   * Test getter
   */
  it('test getter', function () {
    expect( UserStore.getUID() ).toEqual(5);
    expect( UserStore.getIsLogged() ).toEqual(true);
    expect( UserStore.getInfoUser() ).toEqual(wantedResponse);
    expect( UserStore.getInfoMachine() ).toEqual(testInfoMachine);
    expect( UserStore.getMembership() ).toEqual({});
  });

  /*
   * Test submitUpdatedStateToServer
   */
  it('test submitUpdatedStateToServer', function() {
    UserStore.submitUpdatedStateToServer(userUpdateState);
    expect($.ajax).toBeCalledWith({
      headers: {'Content-Type': 'application/json'},
      url: '/api/users/5',
      type: 'PUT',
      data: JSON.stringify(userUpdateFullState),
      success: jasmine.any(Function),
      error: jasmine.any(Function)
    });
  });

  /*
   * Test formatUserStateToSendToServer
   */
  it('test formatUserStateToSendToServer', function() {
    expect( UserStore.formatUserStateToSendToServer(userUpdateState) ).toEqual(userUpdateFullState);
  });

  /*
   * Test cleanState
   * TODO: change toEqual to something which match
   */
  it('test cleanState', function() {
    UserStore.cleanState();
    expect(UserStore._state).toEqual(emptyState);
  });

});

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
      Id: 1,
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
    rawInfoMembership: []
  };

  var testInfoMachine = UserStore._state['rawInfoMachine'];
  var uidTest = 3;

  var wantedResponse = {
    FirstName: "Regular",
    LastName: "User",
    Username: "user",
    Email: "user@example.com",
    InvoiceAddr: 0,
    ShipAddr: 0
  };

  var emptyState = {
    userId : 0,
    isLogged : false,
    rawInfoUser : {},
    rawInfoMachine : [],
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
  it('call into $.ajax with no param or uid', function() {
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

  /*
   *
   * TEST API CALLS
   * with json param
   *
   */
  it('call into $.ajax with json parameters', function() {


    //UserStore.submitStateToServer(userStateTest);

    /*
     * Test submitLoginFormToServer
     */
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
   * Test cleanState
   * TODO: change toEqual to something which match
   */
  it('test cleanState', function() {
    UserStore.cleanState();
    expect(UserStore._state).toEqual(emptyState);
  });
});

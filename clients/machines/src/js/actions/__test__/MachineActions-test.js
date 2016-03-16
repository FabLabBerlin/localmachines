jest.dontMock('../../actionTypes');
jest.dontMock('../ApiActions');
jest.dontMock('../../getters');
jest.dontMock('../LocationActions');
jest.dontMock('../MachineActions');
jest.dontMock('nuclear-js');
jest.dontMock('../../reactor');
jest.dontMock('../../stores/LocationStore');
jest.dontMock('../../stores/LoginStore');
jest.mock('jquery');

var $ = require('jquery');
var actionTypes = require('../../actionTypes');
var LocationActions = require('../LocationActions');
var MachineActions = require('../MachineActions');
var reactor = require('../../reactor');
var LocationStore = require('../../stores/LocationStore');
var LoginStore = require('../../stores/LoginStore');


reactor.registerStores({
  locationStore: LocationStore,
  loginStore: LoginStore
});

LocationActions.setLocationId(1);


describe('MachineActions', function() {
  describe('endActivation', function() {
    it('should POST /api/activations/:aid/close', function() {
      MachineActions.endActivation(2);
      expect($.ajax).toBeCalledWith({
        url: '/api/activations/2/close',
        data: jasmine.any(Object),
        method: 'POST',
        success: jasmine.any(Function),
        error: jasmine.any(Function)
      });
    });
  });

  describe('startActivation', function() {
    it('should POST /api/activations', function() {
      MachineActions.startActivation(17);
      expect($.ajax).toBeCalledWith({
        url: '/api/activations?location=1',
        data: {
          mid: 17
        },
        dataType: 'json',
        type: 'POST',
        success: jasmine.any(Function),
        error: jasmine.any(Function)
      });
    });
  });

  describe('pollActivations', function() {
    it('GETs /api/users/:uid/dashboard', function() {
      var data = {
        UserId: 11
      };
      reactor.dispatch(actionTypes.SUCCESS_LOGIN, { data });
      MachineActions.pollDashboard(1);
      expect($.ajax).toBeCalledWith({
        url: '/api/users/11/dashboard?location=1',
        dataType: 'json',
        type: 'GET',
        cache: false,
        success: jasmine.any(Function),
        error: jasmine.any(Function)
      });
    });
  });
});

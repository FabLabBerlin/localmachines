jest.dontMock('../../actionTypes');
jest.dontMock('../ApiActions');
jest.dontMock('../../getters');
jest.dontMock('../MachineActions');
jest.dontMock('nuclear-js');
jest.dontMock('../../reactor');
jest.dontMock('../../stores/LoginStore');
jest.mock('jquery');

var $ = require('jquery');
var actionTypes = require('../../actionTypes');
var MachineActions = require('../MachineActions');
var reactor = require('../../reactor');
var LoginStore = require('../../stores/LoginStore');


reactor.registerStores({
  loginStore: LoginStore
});


describe('MachineActions', function() {
  describe('endActivation', function() {
    it('should PUT /api/activations/:aid', function() {
      MachineActions.endActivation(2);
      expect($.ajax).toBeCalledWith({
        url: '/api/activations/2',
        data: jasmine.any(Object),
        method: 'PUT',
        success: jasmine.any(Function),
        error: jasmine.any(Function)
      });
    });
  });

  describe('startActivation', function() {
    it('should POST /api/activations', function() {
      MachineActions.startActivation(17);
      expect($.ajax).toBeCalledWith({
        url: '/api/activations',
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

  describe('forceTurnOffMachine', function() {
    it('should POST /api/machines/:mid/turn_off', function() {
      MachineActions.forceTurnOffMachine(17, 2);
      expect($.ajax).toBeCalledWith({
        url: '/api/machines/17/turn_off',
        type: 'POST',
        success: jasmine.any(Function),
        error: jasmine.any(Function)
      });
    });
  });

  describe('forceTurnOnMachine', function() {
    it('should POST /api/machines/:mid/turn_on', function() {
      MachineActions.forceTurnOnMachine(17);
      expect($.ajax).toBeCalledWith({
        url: '/api/machines/17/turn_on',
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
      MachineActions.pollDashboard();
      expect($.ajax).toBeCalledWith({
        url: '/api/users/11/dashboard',
        dataType: 'json',
        type: 'GET',
        cache: false,
        success: jasmine.any(Function),
        error: jasmine.any(Function)
      });
    });
  });
});

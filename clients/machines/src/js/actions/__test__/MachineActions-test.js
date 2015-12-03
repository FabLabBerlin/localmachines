jest.dontMock('../../actionTypes');
jest.dontMock('../ApiActions');
jest.dontMock('../MachineActions');
jest.dontMock('nuclear-js');
jest.mock('jquery');
jest.mock('../../reactor');

var $ = require('jquery');
var actionTypes = require('../../actionTypes');
var MachineActions = require('../MachineActions');
var reactor = require('../../reactor');

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

  describe('adminTurnOffMachine', function() {
    it('should POST /api/machines/:mid/turn_off', function() {
      MachineActions.adminTurnOffMachine(17, 2);
      expect($.ajax).toBeCalledWith({
        url: '/api/machines/17/turn_off',
        data: jasmine.any(Object),
        dataType: 'json',
        type: 'POST',
        success: jasmine.any(Function),
        error: jasmine.any(Function)
      });
    });
  });

  describe('adminTurnOnMachine', function() {
    it('should POST /api/machines/:mid/turn_on', function() {
      MachineActions.adminTurnOnMachine(17);
      expect($.ajax).toBeCalledWith({
        url: '/api/machines/17/turn_on',
        data: jasmine.any(Object),
        dataType: 'json',
        type: 'POST',
        success: jasmine.any(Function),
        error: jasmine.any(Function)
      });
    });
  });

  describe('clearState', function() {
    it('triggers action MACHINE_STORE_CLEAR_STATE', function() {
      MachineActions.clearState();
      expect(reactor.dispatch).toBeCalledWith(actionTypes.MACHINE_STORE_CLEAR_STATE);
    });
  });

  describe('pollActivations', function() {
    it('GETs /api/users/active', function() {
      MachineActions.pollDashboard();
      expect($.ajax).toBeCalledWith({
        url: '/api/activations/active',
        dataType: 'json',
        type: 'GET',
        cache: false,
        success: jasmine.any(Function),
        error: jasmine.any(Function)
      });
    });
  });
});

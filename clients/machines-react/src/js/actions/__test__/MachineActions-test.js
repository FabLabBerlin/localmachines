jest.mock('jquery');
jest.mock('../../flux');
var $ = require('jquery');
var actionTypes = require('../../actionTypes');
var MachineActions = require('../MachineActions');
var Flux = require('../../flux');

describe('MachineActions', function() {
  describe('fetchData', function() {
    it('should GET /api/users/:uid', function() {
      MachineActions.fetchData(123);
      expect($.ajax).toBeCalledWith({
        url: '/api/users/123',
        dataType: 'json',
        type: 'GET',
        cache: false,
        success: jasmine.any(Function),
        error: jasmine.any(Function)
      });
    });
  });

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
      expect(Flux.dispatch).toBeCalledWith(actionTypes.MACHINE_STORE_CLEAR_STATE);
    });
  });

  describe('pollActivations', function() {
    it('GETs /api/activations/active', function() {
      MachineActions.pollActivations();
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

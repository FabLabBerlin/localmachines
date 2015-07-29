jest.mock('jquery');
jest.dontMock('../MachineStore.js');

describe('MachineStore', function() {
  var $ = require('jquery');
  var MachineStore = require('../MachineStore');

  function apiGet(url) {
    return {
      url: url || jasmine.any(String),
      dataType: 'json',
      type: 'GET',
      cache: false,
      success: jasmine.any(Function),
      error: jasmine.any(Function)
    };
  }

  function apiPost(url, data) {
    return {
      url: url || jasmine.any(String),
      dataType: 'json',
      type: 'POST',
      data: data || jasmine.any(Object),
      success: jasmine.any(Function),
      error: jasmine.any(Function)
    };
  }

  function apiPut(url, data) {
    return {
      url: url || jasmine.any(String),
      method: 'PUT',
      data: data || jasmine.any(Object),
      success: jasmine.any(Function),
      error: jasmine.any(Function)
    };
  }

  describe('apiGetUserInfoLogin', function() {
    it('should GET /api/users/:uid', function() {
      MachineStore.apiGetUserInfoLogin(123);
      expect($.ajax).toBeCalledWith(apiGet('/api/users/123'));
    });
  });

  describe('apiGetUserMachines', function() {
    it('should GET /api/users/:uid/machines', function() {
      MachineStore.apiGetUserMachines(123);
      expect($.ajax).toBeCalledWith(apiGet('/api/users/123/machines'));
    });
  });

  describe('apiGetActivationActive', function() {
    it('should GET /api/activations/active', function() {
      MachineStore.apiGetActivationActive();
      expect($.ajax).toBeCalledWith(apiGet('/api/activations/active'));
    });
  });

  describe('apiPostActivation', function() {
    it('should POST /api/activations', function() {
      MachineStore.apiPostActivation(123);
      expect($.ajax).toBeCalledWith(apiPost('/api/activations'));
    });
  });

  describe('apiPutActivation', function() {
    it('should PUT /api/activations/:aid', function() {
      MachineStore.apiPutActivation(123);
      expect($.ajax).toBeCalledWith(apiPut('/api/activations/123'));
    });
  });

  describe('apiPostSwitchMachine', function() {
    describe('on', function() {
      it('should POST /api/machines/:mid/turn_on', function() {
        MachineStore.apiPostSwitchMachine(123, 'on');
        expect($.ajax).toBeCalledWith(apiPost('/api/machines/123/turn_on'));
      });
    });
  });
});

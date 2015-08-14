jest.dontMock('../../flux');
import actionTypes from '../../actionTypes';
import Flux from '../../flux';
import getters from '../../getters';
jest.mock('jquery');
jest.dontMock('../MachineStore.js');

function activation() {
  return [
    {
      FirstName: 'Regular',
      Id: 227,
      LastName: 'Admin',
      MachineId: 2,
      TimeTotal: 1,
      UserId: 2
    }
  ];
}

function machines() {
  return [
    {
      Id: 1,
      Name: 'Printer5000'
    },
    {
      Id: 2,
      Name: 'Form1'
    }
  ];
}

function user() {
  return {
    FirstName: 'Regular',
    Id: 2,
    LastName: 'Admin',
    UserRole: 'admin'
  };
}

describe('MachineStore', function() {
  var $ = require('jquery');
  var LoginStore = require('../LoginStore');
  var MachineStore = require('../MachineStore');

  Flux.registerStores({
    loginStore: LoginStore,
    machineStore: MachineStore
  });

  describe('SET_ACTIVATION_INFO', function() {
    it('does changes visible via getActivationInfo', function() {
      var activationInfo = activation();
      Flux.dispatch(actionTypes.SET_ACTIVATION_INFO, { activationInfo });
      var actual = Flux.evaluateToJS(getters.getActivationInfo);
      expect(actual).toEqual(activation());
    });
  });

  describe('SET_USER_INFO', function() {
    it('does changes visible via getUserInfo', function() {
      var userInfo = user();
      Flux.dispatch(actionTypes.SET_USER_INFO, { userInfo });
      var actual = Flux.evaluateToJS(getters.getUserInfo);
      expect(actual).toEqual(user());
    });
  });

  describe('SET_MACHINE_INFO', function() {
    it('does changes visible via getMachineInfo', function() {
      var machineInfo = machines();
      Flux.dispatch(actionTypes.SET_MACHINE_INFO, { machineInfo });
      var actual = Flux.evaluateToJS(getters.getMachineInfo);
      expect(actual).toEqual(machines());
    });
  });

  describe('MACHINE_STORE_CLEAR_STATE', function() {
    it('clears the state', function() {
      Flux.dispatch(actionTypes.MACHINE_STORE_CLEAR_STATE);
      expect(Flux.evaluateToJS(getters.getUserInfo)).toEqual({});
      expect(Flux.evaluateToJS(getters.getActivationInfo)).toEqual([]);
      expect(Flux.evaluateToJS(getters.getMachineInfo)).toEqual([]);
    });
  });
});

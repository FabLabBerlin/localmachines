jest.dontMock('nuclear-js');
jest.dontMock('../../actionTypes');
jest.dontMock('../../getters');
jest.dontMock('lodash');
jest.dontMock('../LoginStore.js');
jest.dontMock('../MachineStore.js');
jest.dontMock('../../reactor');
jest.mock('jquery');

var actionTypes = require('../../actionTypes');
var getters = require('../../getters');
var reactor = require('../../reactor');


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

  reactor.registerStores({
    loginStore: LoginStore,
    machineStore: MachineStore
  });

  describe('SET_ACTIVATION_INFO', function() {
    it('does changes visible via getActivationInfo', function() {
      var activationInfo = activation();
      reactor.dispatch(actionTypes.SET_ACTIVATION_INFO, { activationInfo });
      var actual = reactor.evaluateToJS(getters.getActivationInfo);
      expect(actual).toEqual(activation());
    });
  });

  describe('SET_USER_INFO', function() {
    it('does changes visible via getUserInfo', function() {
      var userInfo = user();
      reactor.dispatch(actionTypes.SET_USER_INFO, { userInfo });
      var actual = reactor.evaluateToJS(getters.getUserInfo);
      expect(actual).toEqual(user());
    });
  });

  describe('SET_MACHINE_INFO', function() {
    it('does changes visible via getMachineInfo', function() {
      var machineInfo = machines();
      reactor.dispatch(actionTypes.SET_MACHINE_INFO, { machineInfo });
      var actual = reactor.evaluateToJS(getters.getMachineInfo);
      expect(actual).toEqual(machines());
    });
  });

  describe('MACHINE_STORE_CLEAR_STATE', function() {
    it('clears the state', function() {
      reactor.dispatch(actionTypes.MACHINE_STORE_CLEAR_STATE);
      expect(reactor.evaluateToJS(getters.getUserInfo)).toEqual({});
      expect(reactor.evaluateToJS(getters.getActivationInfo)).toEqual([]);
      expect(reactor.evaluateToJS(getters.getMachineInfo)).toEqual([]);
    });
  });
});

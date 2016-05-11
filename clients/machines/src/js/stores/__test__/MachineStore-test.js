jest.dontMock('nuclear-js');
jest.dontMock('../../actionTypes');
jest.dontMock('../../getters');
jest.dontMock('lodash');
jest.dontMock('../LoginStore.js');
jest.dontMock('../MachineStore.js');
jest.dontMock('../../modules/Machines');
jest.dontMock('../UserStore.js');
jest.dontMock('../../reactor');
jest.mock('jquery');

var actionTypes = require('../../actionTypes');
var getters = require('../../getters');
var Machines = require('../../modules/Machines');
var reactor = require('../../reactor');


function getActivations() {
  return [
    {
      FirstName: 'Regular',
      Id: 227,
      LastName: 'Admin',
      MachineId: 2,
      Quantity: 1,
      UserId: 2
    }
  ];
}

function getMachines() {
  return {
    1: {
      Id: 1,
      Name: 'Printer5000'
    },
    2: {
      Id: 2,
      Name: 'Form1'
    }
  };
}

describe('MachineStore', function() {
  var $ = require('jquery');
  var LoginStore = require('../LoginStore');
  var MachineStore = require('../MachineStore');
  var UserStore = require('../UserStore');

  reactor.registerStores({
    loginStore: LoginStore,
    machineStore: MachineStore,
    UserStore: UserStore
  });

  describe('SET_ACTIVATIONS', function() {
    it('does changes visible via getActivations', function() {
      var activations = getActivations();
      reactor.dispatch(Machines.actionTypes.SET_ACTIVATIONS, { activations });
      var actual = reactor.evaluateToJS(getters.getActivations);
      expect(actual).toEqual(getActivations());
    });
  });

  describe('SET_MACHINES', function() {
    it('does changes visible via getMachines', function() {
      var machines = getMachines();
      reactor.dispatch(Machines.actionTypes.SET_MACHINES, { machines });
      var actual = reactor.evaluateToJS(getters.getMachines);
      expect(actual).toEqual(getMachines());
    });
  });

  describe('MACHINE_STORE_CLEAR_STATE', function() {
    it('clears the state', function() {
      reactor.dispatch(Machines.actionTypes.MACHINE_STORE_CLEAR_STATE);
      expect(reactor.evaluateToJS(getters.getActivations)).toEqual(undefined);
      expect(reactor.evaluateToJS(getters.getMachines)).toEqual({});
    });
  });
});

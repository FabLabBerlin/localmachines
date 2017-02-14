jest.dontMock('nuclear-js');
jest.dontMock('../../actionTypes');
jest.dontMock('../../getters');
jest.dontMock('../../../../getters');
jest.dontMock('lodash');
jest.dontMock('../../../../stores/LoginStore.js');
jest.dontMock('../store');
jest.dontMock('../../../../modules/Machines');
jest.dontMock('../../../../modules/Machines/getters');
jest.dontMock('../../../../stores/UserStore.js');
jest.dontMock('../../../../reactor');
jest.mock('jquery');

import actionTypes from '../../actionTypes';
import getters from '../../../../getters';
import Machines from '../../../../modules/Machines';
import reactor from '../../../../reactor';
import LoginStore from '../../../../stores/LoginStore';
import MachineStore from '../../../../modules/Machines/stores/store';
import UserStore from '../../../../stores/UserStore';


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

  reactor.registerStores({
    loginStore: LoginStore,
    machineStore: MachineStore,
    UserStore: UserStore
  });

  describe('SET_ACTIVATIONS', function() {
    it('does changes visible via getActivations', function() {
      var activations = getActivations();
      reactor.dispatch(Machines.actionTypes.SET_ACTIVATIONS, { activations });
      var actual = reactor.evaluateToJS(Machines.getters.getActivations);
      expect(actual).toEqual(getActivations());
    });
  });

  describe('SET_MACHINES', function() {
    it('does changes visible via getMachines', function() {
      var machines = getMachines();
      reactor.dispatch(Machines.actionTypes.SET_MACHINES, { machines });
      var actual = reactor.evaluateToJS(Machines.getters.getMachines);
      expect(actual).toEqual(getMachines());
    });
  });

  describe('MACHINE_STORE_CLEAR_STATE', function() {
    it('clears the state', function() {
      reactor.dispatch(Machines.actionTypes.MACHINE_STORE_CLEAR_STATE);
      expect(reactor.evaluateToJS(Machines.getters.getActivations)).toEqual(undefined);
      expect(reactor.evaluateToJS(Machines.getters.getMachines)).toEqual({});
    });
  });
});

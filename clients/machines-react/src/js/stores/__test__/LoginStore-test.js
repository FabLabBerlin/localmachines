jest.dontMock('../../flux');
import actionTypes from '../../actionTypes';
import Flux from '../../flux';
import getters from '../../getters';
jest.mock('jquery');
jest.mock('toastr');
jest.dontMock('../LoginStore.js');

describe('LoginStore', function() {
  var $ = require('jquery');
  var LoginStore = require('../LoginStore');

  Flux.registerStores({
    loginStore: LoginStore
  });


  describe('SUCCESS_LOGOUT', function() {
    it('sets store into logout state', function() {
      Flux.dispatch(actionTypes.SUCCESS_LOGOUT);
      const isLogged = Flux.evaluateToJS(getters.getIsLogged);
      expect(isLogged).toBe(false);
    });
  });

  describe('SUCCESS_LOGIN', function() {
    it('sets store into logged in state', function() {
      var data = {
        UserId: 123
      };
      Flux.dispatch(actionTypes.SUCCESS_LOGIN, { data });
      const uid = Flux.evaluateToJS(getters.getUid);
      expect(uid).toBe(123);
    });
  });

  describe('ERROR_LOGIN', function() {
    it('sets firstTry to false', function() {
      Flux.dispatch(actionTypes.ERROR_LOGIN);
      const firstTry = Flux.evaluateToJS(getters.getFirstTry);
      expect(firstTry).toBe(false);
    });
  });
});

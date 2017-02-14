jest.dontMock('nuclear-js');
jest.dontMock('../../getters');
jest.dontMock('../LoginStore.js');
jest.dontMock('../../reactor');
jest.mock('jquery');
jest.mock('toastr');

var $ = require('jquery');
import actionTypes from '../../actionTypes';
import getters from '../../getters';
import LoginStore from '../LoginStore';
import reactor from '../../reactor';


describe('LoginStore', function() {
  reactor.registerStores({
    loginStore: LoginStore
  });

  describe('SUCCESS_LOGOUT', function() {
    it('sets store into logout state', function() {
      reactor.dispatch(actionTypes.SUCCESS_LOGOUT);
      const isLogged = reactor.evaluateToJS(getters.getIsLogged);
      expect(isLogged).toBe(false);
    });
  });

  describe('SUCCESS_LOGIN', function() {
    it('sets store into logged in state', function() {
      var data = {
        UserId: 123
      };
      reactor.dispatch(actionTypes.SUCCESS_LOGIN, { data });
      const uid = reactor.evaluateToJS(getters.getUid);
      expect(uid).toBe(123);
    });
  });

  describe('ERROR_LOGIN', function() {
    it('sets firstTry to false', function() {
      reactor.dispatch(actionTypes.ERROR_LOGIN);
      const firstTry = reactor.evaluateToJS(getters.getFirstTry);
      expect(firstTry).toBe(false);
    });
  });
});

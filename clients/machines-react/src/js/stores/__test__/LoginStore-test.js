jest.dontMock('../../flux');
jest.mock('jquery');
jest.mock('toastr');
jest.dontMock('../LoginStore.js');

var $ = require('jquery');
var actionTypes = require('../../actionTypes');
var Flux = require('../../flux');
var getters = require('../../getters');
var LoginStore = require('../LoginStore');


describe('LoginStore', function() {
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

jest.dontMock('../../actionTypes');
jest.dontMock('../../getters');
jest.dontMock('../../modules/Location');
jest.dontMock('../../modules/Location/actions');
jest.dontMock('../../modules/Location/getters');
jest.dontMock('../MachineActions');
jest.dontMock('nuclear-js');
jest.dontMock('../../reactor');
jest.dontMock('../../modules/Location/stores/store');
jest.dontMock('../../stores/LoginStore');
jest.mock('jquery');

var $ = require('jquery');
import actionTypes from '../../actionTypes';
import LocationActions from '../../modules/Location/actions';
import MachineActions from '../MachineActions';
import reactor from '../../reactor';
import LocationStore from '../../modules/Location/stores/store';
import LoginStore from '../../stores/LoginStore';


reactor.registerStores({
  locationStore: LocationStore,
  loginStore: LoginStore
});

LocationActions.setLocationId(1);


describe('MachineActions', function() {
  describe('endActivation', function() {
    it('should POST /api/activations/:aid/close', function() {
      MachineActions.endActivation(2);
      expect($.ajax).toBeCalledWith({
        url: '/api/activations/2/close?location=1',
        data: jasmine.any(Object),
        method: 'POST',
        success: jasmine.any(Function),
        error: jasmine.any(Function)
      });
    });
  });

  describe('startActivation', function() {
    it('should POST /api/activations/start', function() {
      MachineActions.startActivation(17);
      expect($.ajax).toBeCalledWith({
        url: '/api/activations/start?location=1',
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
});

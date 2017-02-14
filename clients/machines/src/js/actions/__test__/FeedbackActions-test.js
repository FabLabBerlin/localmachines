jest.dontMock('../../actionTypes');
jest.dontMock('../FeedbackActions');
jest.dontMock('../../getters');
jest.dontMock('nuclear-js');
jest.dontMock('../../stores/FeedbackStore');
jest.dontMock('../../modules/Machines/stores/store');
jest.dontMock('../../stores/UserStore');

var $ = require('jquery');
import actionTypes from '../../actionTypes';
import FeedbackActions from '../FeedbackActions';
import FeedbackStore from '../../stores/FeedbackStore';
import GlobalActions from '../GlobalActions';
import MachineStore from '../../modules/Machines/stores/store';
import reactor from '../../reactor';
import UserStore from '../../stores/UserStore';


reactor.registerStores({
  feedbackStore: FeedbackStore,
  machineStore: MachineStore,
  userStore: UserStore
});


describe('FeedbackActions', function() {
  describe('reportMachineBroken', function() {
    it('should show the loader and start an ajax request', function() {
      var machineId = 123;
      FeedbackActions.reportMachineBroken({ machineId });
      expect(GlobalActions.showGlobalLoader).toBeCalled();
      expect($.ajax).toBeCalled();
    });
  });

  describe('reportSatisfaction', function() {
    it('should start an ajax request', function() {
      var activationId = 456;
      var satisfaction = 'max';
      FeedbackActions.reportSatisfaction({ activationId, satisfaction });
      expect($.ajax).toBeCalled();
    });
  });

  describe('submit', function() {
    it('should show the loader and start an ajax request', function() {
      FeedbackActions.submit();
      expect(GlobalActions.showGlobalLoader).toBeCalled();
      expect($.ajax).toBeCalled();
    });
  });
});

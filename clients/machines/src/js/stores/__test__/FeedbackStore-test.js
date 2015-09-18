jest.dontMock('nuclear-js');
jest.dontMock('../../actionTypes');
jest.dontMock('../../getters');
jest.dontMock('../../reactor');
jest.dontMock('../../stores/FeedbackStore');
var actionTypes = require('../../actionTypes');
var getters = require('../../getters');
var reactor = require('../../reactor');


describe('FeedbackStore', function() {
  var FeedbackStore = require('../../stores/FeedbackStore');

  reactor.registerStores({
    feedbackStore: FeedbackStore
  });

  it('makes getFeedbackSubject work for billing', function() {
    var key = 'subject-dropdown';
    var value = 'Billing';
    reactor.dispatch(actionTypes.SET_FEEDBACK_PROPERTY, { key, value });
    var subject = reactor.evaluateToJS(getters.getFeedbackSubject);
    expect(subject).toEqual('Billing');
  });

  it('makes getFeedbackSubject work for other text', function() {
    var key = 'subject-dropdown';
    var value = 'Other';
    reactor.dispatch(actionTypes.SET_FEEDBACK_PROPERTY, { key, value });
    key = 'subject-other-text';
    value = 'Helloo';
    reactor.dispatch(actionTypes.SET_FEEDBACK_PROPERTY, { key, value });
    var subject = reactor.evaluateToJS(getters.getFeedbackSubject);
    expect(subject).toEqual('Helloo');
  });
});

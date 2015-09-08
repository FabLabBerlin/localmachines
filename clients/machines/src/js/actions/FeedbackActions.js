var $ = require('jquery');
var actionTypes = require('../actionTypes');
var ApiActions = require('./ApiActions');
var FeedbackStore = require('../stores/FeedbackStore');
var getters = require('../getters');
var reactor = require('../reactor');
var toastr = require('../toastr');

var FeedbackActions = {

  reportSatisfaction({ activationId, satisfaction }) {
    var url = '/api/activations/' + activationId + '/feedback';
    $.ajax({
      url: url,
      dataType: 'json',
      type: 'POST',
      data: {
        satisfaction: satisfaction
      },
      success: function() {
        toastr.info('Feedback submitted');
      },
      error: function(xhr, status, err) {
        console.error(url, status, err.toString());
      }
    });
  },

  setFeedbackProperty({ key, value }) {
    reactor.dispatch(actionTypes.SET_FEEDBACK_PROPERTY, { key, value });
  },

  submit() {
    var userInfo = reactor.evaluateToJS(getters.getUserInfo);
    ApiActions.showGlobalLoader();
    $.ajax({
      url: '/api/feedback',
      dataType: 'json',
      type: 'POST',
      data: {
        subject: reactor.evaluateToJS(getters.getFeedbackSubject),
        message: reactor.evaluateToJS(getters.getFeedbackMessage),
        email: userInfo.Email,
        name: userInfo.FirstName + ' ' + userInfo.LastName
      },
      success: function() {
        reactor.dispatch(actionTypes.RESET_FEEDBACK_FORM);
        ApiActions.hideGlobalLoader();
      },
      error: function(xhr, status, err) {
        toastr.error('Error submitting feedback.  Please try again later.');
        console.error('/feedback', status, err.toString());
        ApiActions.hideGlobalLoader();
      }
    });
  }

};

export default FeedbackActions;

var $ = require('jquery');
var actionTypes = require('../actionTypes');
var ApiActions = require('./ApiActions');
var FeedbackStore = require('../stores/FeedbackStore');
var getters = require('../getters');
var reactor = require('../reactor');
var toastr = require('../toastr');

var FeedbackActions = {

  reportMachineBroken({ machineId }) {
    var machine = getters.getMachine(machineId);
    var userInfo = reactor.evaluateToJS(getters.getUserInfo);
    var fullName = userInfo.FirstName + ' ' + userInfo.LastName;
    ApiActions.showGlobalLoader();
    $.ajax({
      url: '/api/feedback',
      dataType: 'json',
      type: 'POST',
      data: {
        subject: 'Machine Broken',
        message: 'Hi friends,\n\nThe following machine seems broken: ' + machine.Name + '\n\nEasyLab on behalf of ' + fullName,
        email: userInfo.Email,
        name: fullName
      },
      success: function() {
        toastr.info('Thank you for the report 😀 We will have a look at it asap.');
        ApiActions.hideGlobalLoader();
      },
      error: function(xhr, status, err) {
        toastr.error('Error submitting report.  Please try again later.');
        console.error('/feedback', status, err.toString());
        ApiActions.hideGlobalLoader();
      }
    });
  },

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

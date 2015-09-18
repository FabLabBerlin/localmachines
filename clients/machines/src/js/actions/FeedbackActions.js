var $ = require('jquery');
var actionTypes = require('../actionTypes');
var ApiActions = require('./ApiActions');
var FeedbackStore = require('../stores/FeedbackStore');
var getters = require('../getters');
var LoginActions = require('../actions/LoginActions');
var reactor = require('../reactor');
var toastr = require('../toastr');

var FeedbackActions = {

  reportMachineBroken({ machineId }) {
    LoginActions.keepAlive();
    const machinesById = reactor.evaluateToJS(getters.getMachinesById);
    const machine = machinesById[machineId] || {};
    var userInfo = reactor.evaluateToJS(getters.getUserInfo);
    var fullName = userInfo.FirstName + ' ' + userInfo.LastName;
    ApiActions.showGlobalLoader();
    $.ajax({
      url: '/api/machines/' + machine.Id + '/report_broken',
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
    LoginActions.keepAlive();
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
    LoginActions.keepAlive();
    reactor.dispatch(actionTypes.SET_FEEDBACK_PROPERTY, { key, value });
  },

  submit() {
    LoginActions.keepAlive();
    var userInfo = reactor.evaluateToJS(getters.getUserInfo);
    var subject = reactor.evaluateToJS(getters.getFeedbackSubject);
    var message = reactor.evaluateToJS(getters.getFeedbackMessage);
    ApiActions.showGlobalLoader();
    if (subject && message) {
      $.ajax({
        url: '/api/feedback',
        dataType: 'json',
        type: 'POST',
        data: {
          subject: subject,
          message: message,
          email: userInfo.Email,
          name: userInfo.FirstName + ' ' + userInfo.LastName
        },
        success: function() {
          reactor.dispatch(actionTypes.RESET_FEEDBACK_FORM);
          ApiActions.hideGlobalLoader();
          toastr.info('Thank you for your feedback 😀');
        },
        error: function(xhr, status, err) {
          toastr.error('Error submitting feedback.  Please try again later.');
          console.error('/feedback', status, err.toString());
          ApiActions.hideGlobalLoader();
        }
      });
    } else {
      toastr.error('Please fill out all fields');
    }
  }

};

export default FeedbackActions;

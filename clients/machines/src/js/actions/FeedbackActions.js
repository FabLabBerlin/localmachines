var $ = require('jquery');
import actionTypes from '../actionTypes';
import FeedbackStore from '../stores/FeedbackStore';
import getters from '../getters';
import GlobalActions from './GlobalActions';
import LocationGetters from '../modules/Location/getters';
import LoginActions from '../actions/LoginActions';
import Machines from '../modules/Machines';
import reactor from '../reactor';
import toastr from '../toastr';

var FeedbackActions = {

  reportMachineBroken({ machineId, text }) {
    LoginActions.keepAlive();
    const machinesById = reactor.evaluateToJS(Machines.getters.getMachinesById);
    const machine = machinesById[machineId] || {};
    var user = reactor.evaluateToJS(getters.getUser);
    var fullName = user.FirstName + ' ' + user.LastName;
    GlobalActions.showGlobalLoader();

    $.ajax({
      url: '/api/machines/' + machine.Id + '/report_broken',
      dataType: 'json',
      type: 'POST',
      data: {
        subject: 'Machine Broken',
        message: 'Hi friends,\n\nThe following machine seems broken: ' + machine.Name + '\n\nThe user wrote:' + text + '\n\nEASY LAB on behalf of ' + fullName,
        text: text,
        email: user.Email,
        name: fullName
      },
      success() {
        toastr.info('Thank you for the report 😀 We will have a look at it asap.');
        GlobalActions.hideGlobalLoader();
      },
      error(xhr, status, err) {
        toastr.error('Error submitting report.  Please try again later.');
        console.error('/feedback', status, err.toString());
        GlobalActions.hideGlobalLoader();
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
      success() {
        toastr.info('Feedback submitted');
      },
      error(xhr, status, err) {
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
    var user = reactor.evaluateToJS(getters.getUser);
    var subject = reactor.evaluateToJS(getters.getFeedbackSubject);
    var message = reactor.evaluateToJS(getters.getFeedbackMessage);
    var locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    if (subject && message) {
      GlobalActions.showGlobalLoader();
      $.ajax({
        url: '/api/feedback?location=' + locationId,
        dataType: 'json',
        type: 'POST',
        data: {
          subject: subject,
          message: message,
          email: user.Email,
          name: user.FirstName + ' ' + user.LastName
        },
        success() {
          reactor.dispatch(actionTypes.RESET_FEEDBACK_FORM);
          GlobalActions.hideGlobalLoader();
          toastr.info('Thank you for your feedback 😀');
        },
        error(xhr, status, err) {
          toastr.error('Error submitting feedback.  Please try again later.');
          console.error('/feedback', status, err.toString());
          GlobalActions.hideGlobalLoader();
        }
      });
    } else {
      toastr.error('Please fill out all fields');
    }
  }

};

export default FeedbackActions;

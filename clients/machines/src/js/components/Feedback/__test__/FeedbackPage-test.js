jest.dontMock('lodash');
jest.dontMock('moment');
jest.dontMock('nuclear-js');
jest.dontMock('react');
jest.dontMock('../../../actionTypes');
jest.dontMock('../FeedbackPage');
jest.dontMock('../../../getters');
jest.dontMock('../../../stores/FeedbackStore');
jest.dontMock('../../../stores/LoginStore');


import actionTypes from '../../../actionTypes';
import FeedbackActions from '../../../actions/FeedbackActions';
import FeedbackPage from '../FeedbackPage';
import FeedbackStore from '../../../stores/FeedbackStore';
import getters from '../../../getters';
import LoginStore from '../../../stores/LoginStore';
import Nuclear from 'nuclear-js';
import React from 'react';
import reactor from '../../../reactor';
var toImmutable = Nuclear.toImmutable;


reactor.registerStores({
  loginStore: LoginStore,
  feedbackStore: FeedbackStore
});


describe('FeedbackPage', function() {
  describe('render', function() {
    it('can submit Technical feedback', function() {
      import TestUtils from 'react-addons-test-utils';
      var feedbackPage = TestUtils.renderIntoDocument(
        <FeedbackPage />
      );
      var subject = TestUtils.findRenderedDOMComponentWithTag(feedbackPage, 'select');
      var message = TestUtils.findRenderedDOMComponentWithTag(feedbackPage, 'textarea');
      var submitButton = TestUtils.findRenderedDOMComponentWithTag(feedbackPage, 'button');
      TestUtils.Simulate.change(subject, {
        target: {
          id: 'subject-dropdown',
          value: 'Technical'
        }
      });
      var key = 'subject-dropdown';
      var value = 'Technical';
      expect(FeedbackActions.setFeedbackProperty).toBeCalledWith({ key, value });
      TestUtils.Simulate.change(message, {
        target: {
          id: 'message',
          value: 'Error123'
        }
      });

      key = 'message';
      value = 'Error123';
      expect(FeedbackActions.setFeedbackProperty).toBeCalledWith({ key, value });
      TestUtils.Simulate.click(submitButton);
      expect(FeedbackActions.submit).toBeCalled();
    });
  });
});

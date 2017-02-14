jest.dontMock('nuclear-js');
jest.dontMock('react');
jest.dontMock('../RepairButton');
jest.dontMock('../../../actionTypes');
jest.dontMock('../../../modules/Machines/actionTypes');


import FeedbackDialogs from '../../Feedback/FeedbackDialogs';
import React from 'react';
import RepairButton from '../RepairButton';


describe('RepairButton', function() {
	describe('render', function() {
    it('renders the button and makes it clickable', function() {
      import TestUtils from 'react-addons-test-utils';
      var repairButton = TestUtils.renderIntoDocument(
        <RepairButton/>
      );
      var button = TestUtils.findRenderedDOMComponentWithTag(
        repairButton, 'a');
      TestUtils.Simulate.click(button);
      expect(FeedbackDialogs.machineIssue).toBeCalled();
    });
  });
});

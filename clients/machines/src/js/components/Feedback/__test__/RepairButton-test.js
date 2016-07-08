jest.dontMock('nuclear-js');
jest.dontMock('react');
jest.dontMock('../RepairButton');
jest.dontMock('../../../actionTypes');
jest.dontMock('../../../modules/Machines/actionTypes');


var FeedbackDialogs = require('../../Feedback/FeedbackDialogs');
var React = require('react');
var RepairButton = require('../RepairButton');


describe('RepairButton', function() {
	describe('render', function() {
    it('renders the button and makes it clickable', function() {
      var TestUtils = require('react-addons-test-utils');
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

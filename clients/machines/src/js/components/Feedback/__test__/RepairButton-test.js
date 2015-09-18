jest.dontMock('../RepairButton');


var FeedbackDialogs = require('../../Feedback/FeedbackDialogs');
var RepairButton = require('../RepairButton');


describe('RepairButton', function() {
	describe('render', function() {
    it('renders the button and makes it clickable', function() {
      var React = require('react/addons');
      var TestUtils = React.addons.TestUtils;
      var repairButton = TestUtils.renderIntoDocument(
        <RepairButton/>
      );
      var button = TestUtils.findRenderedDOMComponentWithTag(repairButton, 'a');
      TestUtils.Simulate.click(button);
      expect(FeedbackDialogs.machineIssue).toBeCalled();
    });
  });
});

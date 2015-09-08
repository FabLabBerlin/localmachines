var FeedbackDialogs = require('./FeedbackDialogs');
var React = require('react');


var RepairButton = React.createClass({
  handleClick() {
    FeedbackDialogs.machineIssue(this.props.machineId);
  },

  render() {
    return (
      <button className="machine-report-problem btn btn-sm btn-secondary"
              onClick={this.handleClick}>
        <i className="fa fa-wrench"/>
      </button>
    );
  }
});

export default RepairButton;

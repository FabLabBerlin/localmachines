var FeedbackDialogs = require('./FeedbackDialogs');
var React = require('react');


var RepairButton = React.createClass({
  handleClick() {
    FeedbackDialogs.machineIssue(this.props.machineId);
  },

  render() {
    return (
      <div className="machine-repair-switch row">
        <a
          className="danger" 
          href="#" 
          onClick={this.handleClick}>
          <i className="fa fa-wrench"></i>
        </a>Request Repair
      </div>
    );
  }
});

export default RepairButton;

var FeedbackDialogs = require('./FeedbackDialogs');
var React = require('react');


var RepairButton = React.createClass({
  handleClick() {
    FeedbackDialogs.machineIssue(this.props.machineId);
  },

  render() {
    return (
      <div className="machine-service-repair">
        <a
          className="danger" 
          href="#" 
          onClick={this.handleClick}>
          <i className="fa fa-wrench"></i>
        </a>Report a machine fault
      </div>
    );
  }
});

export default RepairButton;

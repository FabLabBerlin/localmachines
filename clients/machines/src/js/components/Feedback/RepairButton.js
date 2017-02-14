import FeedbackDialogs from './FeedbackDialogs';
import React from 'react';


var RepairButton = React.createClass({
  handleClick() {
    FeedbackDialogs.machineIssue(this.props.machineId);
  },

  render() {
    return (
      <div className="machine-service-repair">
        <a
          className="danger" 
          href="javascript:void(0)" 
          onClick={this.handleClick}>
          <i className="fa fa-wrench"></i>
        </a>Report a machine failure
      </div>
    );
  }
});

export default RepairButton;

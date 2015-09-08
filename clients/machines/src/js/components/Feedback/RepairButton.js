var React = require('react');


var RepairButton = React.createClass({
  handleClick() {
    console.log('Click!');
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

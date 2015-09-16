var React = require('react');


var ForceSwitchOff = React.createClass({

  handleForceSwitchOff() {
    this.props.force('off');
  },

  render() {

    return (
      <div className="machine-force-switch-off">
        <a href="#" onClick={this.handleForceSwitchOff}>
          <i className="fa fa-power-off"></i>
        </a> Force Switch Off
      </div>
    );

  }

});

export default ForceSwitchOff;

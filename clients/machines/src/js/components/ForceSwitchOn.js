var React = require('react');


var ForceSwitchOn = React.createClass({

  handleForceSwitchOn() {
    this.props.force('on');
  },

  render() {

    return (
      <div className="machine-force-switch-on">
        <a href="#" onClick={this.handleForceSwitchOn}>
          <i className="fa fa-power-off"></i>
        </a> Force Switch On
      </div>
    );

  }

});

export default ForceSwitchOn;

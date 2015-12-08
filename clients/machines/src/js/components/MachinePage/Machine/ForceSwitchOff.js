var React = require('react');


var ForceSwitchOff = React.createClass({

  handleForceSwitchOff() {
    this.props.force('off');
  },

  render() {

    return (
      <a 
        href="#" 
        onClick={this.handleForceSwitchOff}
        className="force-switch force-switch-off">
        <i className="fa fa-power-off"></i>
      </a>
    );

  }

});

export default ForceSwitchOff;

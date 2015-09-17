var React = require('react');


var ForceSwitchOn = React.createClass({

  handleForceSwitchOn() {
    this.props.force('on');
  },

  render() {

    return (
      <a 
        href="#" 
        onClick={this.handleForceSwitchOn}
        className="force-switch force-switch-on">
        <i className="fa fa-power-off"></i>
      </a>
    );

  }

});

export default ForceSwitchOn;

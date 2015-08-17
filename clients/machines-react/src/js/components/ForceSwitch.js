var React = require('react');


var ForceSwitch = React.createClass({
  /*
  * Force the switch to turn on
  */
  handleForceSwitchOn() {
    this.props.force('on');
  },

  /*
  * Force the switch to trun off
  */
  handleForceSwitchOff() {
    this.props.force('off');
  },

  render() {
    if (this.props.isAdmin) {
      return (
        <div className="pull-right" >
          <span>Force Switch: </span>
          <button
            onClick={this.handleForceSwitchOn}
            className="btn btn-lg btn-primary">On</button>
          <button
            onClick={this.handleForceSwitchOff}
            className="btn btn-lg btn-danger">Off</button>
        </div>
      );
    } else {
      return <div/>;
    }
  }
});

export default ForceSwitch;

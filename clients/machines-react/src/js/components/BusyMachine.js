import React from 'react';
import MachineActions from '../actions/MachineActions';
import Timer from './Timer';

/*
 * Div displayed when a machine is busy
 * can send action to end an activation
 */
var BusyMachine = React.createClass({

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

  /*
   * Send an action to the store to end the activation
   */
  endActivation(event) {
    event.preventDefault();
    this.props.func(this.props.activation.Id);
  },

  /*
   * Render stuff
   * TODO: make commentaries
   */
  render() {
    return(
      <div className="row" >
        <div className="col-xs-6">
          <Timer time={this.props.activation.TimeTotal} />
        </div>
        <div className="col-xs-6" >
          <button 
            className="btn btn-lg btn-danger btn-block"
            onClick={this.endActivation}
            >Stop</button>
          { this.props.isAdmin ? (
            <div className="pull-right" >
              <label>Force Switch</label>
              <button 
                onClick={this.handleForceSwitchOn}
                className="btn btn-lg btn-primary" >On</button>
              <button 
                onClick={this.handleForceSwitchOff}
                className="btn btn-lg btn-danger" >Off</button>
            </div>
          ):('') }
        </div>
      </div>
    );
  }
});

module.exports = BusyMachine;

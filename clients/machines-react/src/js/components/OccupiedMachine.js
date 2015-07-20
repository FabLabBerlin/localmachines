import React from 'react';
import MachineActions from '../actions/MachineActions';
import Timer from './Timer';

var OccupiedMachine = React.createClass({

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
    this.props.func();
  },

  /*
   * Render the occupied div
   * Become a button if the user is an admin
   */
  render() {
    return (
      <div className="container-fluid">
        <div className="col-xs-6" >
          <p>Occupied by</p>
          <label>{this.props.activation.FirstName + this.props.activation.LastName}</label>
          <Timer time={this.props.activation.TimeTotal} />
        </div>

        { this.props.isAdmin ? (
          <div className="col-xs-6" >
            <button
              className="btn btn-lg btn-warning btn-block"
              onClick={this.endActivation}
              >
              Stop
            </button>
            <div className="pull-right" >
              <label>Force Switch</label>
              <button 
                onClick={this.handleForceSwitchOn}
                className="btn btn-lg btn-primary" >On</button>
              <button 
                onClick={this.handleForceSwitchOff}
                className="btn btn-lg btn-danger" >Off</button>
            </div>
          </div>
            ) : (
            <div className="col-xs-6" >
              <div className="indicator indicator-occupied" >occupied</div>
            </div>
            )}
      </div>
    );
  }
});

module.exports = OccupiedMachine;

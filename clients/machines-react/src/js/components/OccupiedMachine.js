import React from 'react';
import ForceSwitch from './ForceSwitch';
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
   * Render Busy div
   * If the machine is occupied by someone else it will be displayed
   * Admin have a button to stop a machine use by an other user
   * Admin have two more button to force switch
   */
  render() {
    return (
      <div className="container-fluid">
        <div className="col-xs-6" >
          <p>Occupied by</p>
          <label>{this.props.activation.FirstName} {this.props.activation.LastName}</label>
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
            <ForceSwitch isAdmin={this.props.isAdmin} force={this.props.force}/>
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

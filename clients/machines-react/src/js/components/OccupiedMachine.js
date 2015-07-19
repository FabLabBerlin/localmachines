import React from 'react';
import MachineActions from '../actions/MachineActions';
import Timer from './Timer';

var OccupiedMachine = React.createClass({

  /*
   * Send an action to the store to end the activation
   */
  endActivation(event) {
    event.preventDefault();
    MachineActions.endActivation(this.props.activation.Id);
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

        <div className="col-xs-6" >
          { this.props.user.Role == 'admin' ? (
            <button
              className="btn btn-lg btn-warning btn-block"
              onClick={this.endActivation}
              >
              Stop
            </button>
            ): (
            <div className="indicator indicator-occupied" >occupied</div>
            )}
        </div>
      </div>
    );
  }
});

module.exports = OccupiedMachine;

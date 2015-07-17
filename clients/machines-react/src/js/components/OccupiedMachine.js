import React from 'react';
import MachineActions from '../actions/MachineActions';

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
      <div>
        <div className="container-fluid">
          {this.props.info.Name}
          <br/>
          {this.props.activation}
        </div>
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
    );
  }
});

module.exports = OccupiedMachine;

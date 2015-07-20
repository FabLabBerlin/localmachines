import React from 'react';
import MachineActions from '../actions/MachineActions';
import Timer from './Timer';

/*
 * Div displayed when a machine is busy
 * can send action to end an activation
 */
var BusyMachine = React.createClass({

  /*
   * Send an action to the store to end the activation
   */
  endActivation(event) {
    event.preventDefault();
    MachineActions.endActivation(this.props.activation.Id);
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
        </div>
      </div>
    );
  }
});

module.exports = BusyMachine;

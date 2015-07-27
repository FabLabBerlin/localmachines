import React from 'react';
import ForceSwitch from './ForceSwitch';
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
    this.props.func();
  },

  /*
   * Render Busy div
   * If the machine is occupied by the user it will be displayed
   * Admin have two more button to force switch
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
          <ForceSwitch isAdmin={this.props.isAdmin} force={this.props.force}/>
        </div>
      </div>
    );
  }
});

module.exports = BusyMachine;

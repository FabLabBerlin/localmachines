import React from 'react';
import MachineActions from '../actions/MachineActions';

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
    var aid = this.props.activation.Id;
    MachineActions.endActivation(aid);
  },

  /*
   * Render stuff
   * TODO: make commentaries
   */
  render() {
    return(
      <div>
        <div className="container-fluid" >
          {this.props.info.Name}
          <br/>
          {this.props.activation}
        </div>
        <button 
          className="btn btn-danger"
          onClick={this.endActivation}
          >stop</button>
      </div>
    );
  }
});

module.exports = BusyMachine;

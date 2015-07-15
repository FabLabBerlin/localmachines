import React from 'react';
import MachineActions from '../actions/MachineActions';

/*
 * Div displayed the machine is free
 * Can activate an activation
 */
var FreeMachine = React.createClass ({

  /*
   * Try to activate the machine
   */
  startActivation() {
    var mid = this.props.info.Id;
    MachineActions.startActivation(mid);
  },

  /*
   * Render stuff
   * TODO: real commentaries
   */
  render() {
    return (
      <div>
        <div className="container-fluid" >
          {this.props.info.Name}
          <br/>
          {this.props.activation}
        </div>
        <button
          className="btn btn-primary"
          onClick={this.startActivation}
          >start </button>
      </div>
    );
  }
});

module.exports = FreeMachine;

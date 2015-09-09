var ForceSwitch = require('./ForceSwitch');
var MachineActions = require('../actions/MachineActions');
var React = require('react');
var RepairButton = require('./Feedback/RepairButton');
var Timer = require('./Timer');


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
    return (
      <div className="row" >
        <div className="col-xs-6">
          <Timer time={this.props.activation.TimeTotal} />
          <RepairButton machineId={this.props.info.Id}/>
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

export default BusyMachine;

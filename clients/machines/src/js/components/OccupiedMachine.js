var ForceSwitch = require('./ForceSwitch');
var MachineActions = require('../actions/MachineActions');
var React = require('react');
var Timer = require('./Timer');


var OccupiedMachine = React.createClass({

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
      <div className="row">
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

export default OccupiedMachine;
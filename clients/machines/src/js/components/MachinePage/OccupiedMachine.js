var ForceSwitch = require('./ForceSwitch');
var getters = require('../getters');
var MachineActions = require('../actions/MachineActions');
var MaintenanceSwitch = require('./MaintenanceSwitch');
var React = require('react');
var RepairButton = require('./Feedback/RepairButton');
var reactor = require('../reactor');
var Timer = require('./Timer');


var OccupiedMachine = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      machineUsers: getters.getMachineUsers
    };
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
    var users = this.state.machineUsers;
    var user = users.get(this.props.activation.UserId) || {};
    return (
      <div className="row">
        <div className="col-xs-6" >
          <p>Occupied by</p>
          <label>{user.FirstName} {user.LastName}</label>
          <Timer time={this.props.activation.TimeTotal} />
          <RepairButton machineId={this.props.info.Id}/>
        </div>

        { this.props.isAdmin ? (
          <div className="col-xs-6" >
            <button
              className="btn btn-lg btn-warning btn-block"
              onClick={this.endActivation}
              >
              Stop
            </button>
            <MaintenanceSwitch machineId={this.props.info.Id}/>
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

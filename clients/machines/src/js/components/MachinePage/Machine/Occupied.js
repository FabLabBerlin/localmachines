var getters = require('../../../getters');
var MachineActions = require('../../../actions/MachineActions');
var Machines = require('../../../modules/Machines');
var React = require('react');
var reactor = require('../../../reactor');
var Timer = require('./Timer');


var OccupiedMachine = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      machineUsers: Machines.getters.getMachineUsers
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
   * Staff have a button to stop a machine use by an other user
   * Staff have two more button to force switch
   */
  render() {
    var users = this.state.machineUsers;
    var user = users.get(this.props.activation.UserId) || {};
    return (
      <div className="machine occupied">
        <div className="row">
          <div className="col-xs-6">
  
            <div className="machine-action-info">
              <div className="machine-info-content">
                <div className="occupied-by-label">
                  Occupied by
                </div>
                <div className="occupied-by-value">
                  {user.FirstName} {user.LastName}
                </div>
                <Timer activation={this.props.activation} />
              </div>
            </div>
  
          </div>
  
          { this.props.isStaff ? (
            <div className="col-xs-6">
  
              <table className="machine-activation-table">
                <tr>
                  <td rowSpan="2">
                    <button
                      className="btn btn-lg btn-warning btn-block"
                      onClick={this.endActivation}>
                      Stop
                    </button>
                  </td>
                </tr>
              </table>
  
            </div>
          ) : (
            <div className="col-xs-6">
              <div className="indicator occupied">Occupied</div>
            </div>
          )}
  
        </div>
      </div>
    );
  }
});

export default OccupiedMachine;

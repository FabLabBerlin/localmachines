var getters = require('../../../getters');
var MachineActions = require('../../../actions/MachineActions');
var React = require('react');
var reactor = require('../../../reactor');
var Timer = require('./Timer');
var ForceSwitchOn = require('./ForceSwitchOn');
var ForceSwitchOff = require('./ForceSwitchOff');


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
  
          { this.props.isAdmin ? (
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
                  <td className="force-button-table-cell">
                    {this.props.isAdmin ? (
                      <ForceSwitchOn force={this.props.force}/>
                    ) : ''}
                  </td>
                </tr>
                <tr>
                  <td className="force-button-table-cell">
                    {this.props.isAdmin ? (
                      <ForceSwitchOff force={this.props.force}/>
                    ) : ''}
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

var getters = require('../../../getters');
var LocationGetters = require('../../../modules/Location/getters');
var MachineActions = require('../../../actions/MachineActions');
var Machines = require('../../../modules/Machines');
var React = require('react');
var reactor = require('../../../reactor');
var Timer = require('./ActivationTimer');


var UnavailableMachine = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      machineUsers: Machines.getters.getMachineUsers
    };
  },

  startActivation() {
    this.props.startActivation();
  },

  /*
   * Send an action to the store to end the activation
   */
  endActivation(event) {
    event.preventDefault();
    this.props.endActivation();
  },

  /*
   * Render Unavailable div
   * If the machine is unavailable due to maintenance works
   */
  render() {
    var users = this.state.machineUsers;
    var user = this.props.activation ? users.get(this.props.activation.UserId, {}) : {};
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    var startStopButton;

    if (this.props.busy) {
      startStopButton = (
        <button 
          className="btn btn-lg btn-default btn-block"
          onClick={this.endActivation}>
          Stop
        </button>
      );
    } else {
      startStopButton = (
        <button 
          className="btn btn-lg btn-default btn-block"
          onClick={this.startActivation}>
          Start
        </button>
      );
    }
    return (
      <div className="machine unavailable">
        <div className="row">
          <div className="col-xs-6">
  
            <div className="machine-action-info">
              <div className="machine-info-content">
                Undergoing maintenance.
                {locationId === 1 ?
                  <span>Check <a target="_blank" href="https://twitter.com/FabLabBLNAI">@FabLabBLNAI</a> for updates.</span>
                  : null}
              </div>
            </div>
  
          </div>
  
          <div className="col-xs-6">

          { this.props.isStaff ? (
            <table className="machine-activation-table">
              <tbody>
                <tr>
                  <td rowSpan="2">
                    {startStopButton}
                  </td>
                </tr>
              </tbody>
            </table>
          ) : (
            <div className="indicator unavailable">
              Unavailable
            </div>
          )}

          </div>
  
        </div>
      </div>
    );
  }
});

export default UnavailableMachine;

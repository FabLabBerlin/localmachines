var getters = require('../../../getters');
var MachineActions = require('../../../actions/MachineActions');
var React = require('react');
var reactor = require('../../../reactor');
var Timer = require('./Timer');
var ForceSwitchOn = require('./ForceSwitchOn');
var ForceSwitchOff = require('./ForceSwitchOff');


var UnavailableMachine = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      machineUsers: getters.getMachineUsers
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
    var user = users.get(this.props.activation.UserId) || {};
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
                Undergoing maintenance. Check <a target="_blank" href="https://twitter.com/FabLabBLNAI">@FabLabBLNAI</a> for updates.
              </div>
            </div>
  
          </div>
  
          <div className="col-xs-6">

          { this.props.isAdmin ? (
            <table className="machine-activation-table">
              <tr>
                <td rowSpan="2">
                  {startStopButton}
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

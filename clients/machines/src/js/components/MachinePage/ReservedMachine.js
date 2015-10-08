var getters = require('../../getters');
var MachineActions = require('../../actions/MachineActions');
var React = require('react');
var reactor = require('../../reactor');
var Timer = require('./Timer');
var ForceSwitchOn = require('../ForceSwitchOn');
var ForceSwitchOff = require('../ForceSwitchOff');


var ReservedMachine = React.createClass({
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
    const isReservator = this.props.reservation.get('UserId') === this.props.user.get('Id');
    return (
      <div className="machine reserved">
        <div className="row">
          <div className="col-xs-6">
  
            <div className="machine-action-info">
              <div className="machine-info-content">
                {isReservator ? (
                  'This machine is reserved by you.'
                ) : (
                  'This machine is reserved.'
                )}
                
              </div>
            </div>
  
          </div>
  
          <div className="col-xs-6">

          { (this.props.isAdmin || isReservator) ? (
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
              Reserved
            </div>
          )}

          </div>
  
        </div>
      </div>
    );
  }
});

export default ReservedMachine;

var getters = require('../../../getters');
var MachineActions = require('../../../actions/MachineActions');
var Machines = require('../../../modules/Machines');
var React = require('react');
var reactor = require('../../../reactor');


var ReservedMachine = React.createClass({
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
    var user = users.get(this.props.reservation.get('UserId')) || {};
    var startStopButton;
    const isReservator = this.props.reservation.get('UserId') === this.props.user.get('Id');
    if (this.props.busy) {
      startStopButton = (
        <button 
          className="btn btn-lg btn-danger btn-block"
          onClick={this.endActivation}>
          Stop
        </button>
      );
    } else {

      if (isReservator) {
        startStopButton = (
          <button 
            className="btn btn-lg btn-reservation btn-block"
            onClick={this.startActivation}>
            Start
          </button>
        );
      } else {
        startStopButton = (
          <button 
            className="btn btn-lg btn-reservation btn-block"
            onClick={this.startActivation}>
            Start
          </button>
        );
      }
    }
    
    const reservedClassName = 'machine reserved' + 
      (isReservator ? ' reservator' : '');
    return (
      <div className={reservedClassName}>
        <div className="row">
          <div className="col-xs-6">
  
            <div className="machine-action-info">
              <div className="machine-info-content">
                <div className="reserved-by-label">
                  Reserved by
                </div>
                <div className="reserved-by-value">
                {isReservator ? 
                  ('You') : 
                  (user.FirstName + ' ' + user.LastName)
                }
                </div>
              </div>
            </div>
  
          </div>
  
          <div className="col-xs-6">

          { (isReservator && !this.props.isStaff) ? ( {startStopButton} ) : '' }

          { (this.props.isStaff) ? (
            <table className="machine-activation-table">
              <tbody>
                <tr>
                  <td rowSpan="2">
                    {startStopButton}
                  </td>
                </tr>
              </tbody>
            </table>
          ) : ''}

          {(!this.props.isStaff && !isReservator) ? (
            <div className="indicator reserved">
              Reserved
            </div>
          ) : ''}

          </div>
  
        </div>
      </div>
    );
  }
});

export default ReservedMachine;

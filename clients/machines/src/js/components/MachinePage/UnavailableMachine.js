var getters = require('../../getters');
var MachineActions = require('../../actions/MachineActions');
var React = require('react');
var reactor = require('../../reactor');
var Timer = require('./Timer');
var ForceSwitchOn = require('../ForceSwitchOn');
var ForceSwitchOff = require('../ForceSwitchOff');


var UnavailableMachine = React.createClass({
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
   * Render Unavailable div
   * If the machine is unavailable due to maintenance works
   */
  render() {
    var users = this.state.machineUsers;
    var user = users.get(this.props.activation.UserId) || {};
    return (
      <div className="machine unavailable">
        <div className="row">
          <div className="col-xs-6">
            
            <div className="machine-options-toggle" />
  
            <div className="machine-action-info">
              <div className="machine-info-content">
              </div>
            </div>
  
          </div>
  
          { this.props.isAdmin ? (
            <div className="col-xs-6">
  
              <table className="machine-activation-table">
                <tr>
                  <td rowSpan="2">
                    <div className="machine-action-info">
                      <div className="machine-info-content machine-info-unavailable">
                        Machine is unavailable due to maintenance works.
                        Check <a href="https://twitter.com/FabLabBLNAI">@FabLabBLNAI</a> to see
                        when it works again.
                      </div>
                    </div>
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
              <div className="indicator indicator-unavailable">Unavailable</div>
            </div>
          )}
  
        </div>
      </div>
    );
  }
});

export default UnavailableMachine;

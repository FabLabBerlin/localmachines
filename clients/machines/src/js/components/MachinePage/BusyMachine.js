var React = require('react');
var RepairButton = require('../Feedback/RepairButton');
var Timer = require('./Timer');
var ForceSwitch = require('./ForceSwitch');
var ForceSwitchOn = require('../ForceSwitchOn');
var ForceSwitchOff = require('../ForceSwitchOff');

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
      <div className="machine used">
        <table className="row">
  
          <th>
  
            <div>
              <div className="machine-info-content row">
                <div className="machine-time-label">
                  Usage time
                </div>
                <Timer time={this.props.activation.TimeTotal} />
              </div>
            </div>
  
          </th>
  
          <td className="col-xs-12">
  
          { this.props.isAdmin ? (
            
            <table className="machine-activation-table">
              <tr>
                <td rowSpan="2">
                  <button
                    className="btn-stop btn btn-lg btn-danger btn-block"
                    onClick={this.endActivation}>
                    Stop
                  </button>
                </td>
                <td className="force-button-table-cell">
                  <ForceSwitchOn force={this.props.force}/>
                </td>
              </tr>
              <tr>
                <td className="force-button-table-cell">
                  <ForceSwitchOff force={this.props.force}/>
                </td>
              </tr>
            </table>
            
          ) : (
  
            <button
              className="btn-stop btn btn-lg btn-danger btn-block"
              onClick={this.endActivation}>
              Stop
            </button>
  
          )}
  
          </td>
  
        </table>
      </div>
    );
  }
});

export default BusyMachine;

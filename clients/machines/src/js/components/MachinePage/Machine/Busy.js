var React = require('react');
var RepairButton = require('../../Feedback/RepairButton');
var Timer = require('./Timer');
var ForceSwitchOn = require('./ForceSwitchOn');
var ForceSwitchOff = require('./ForceSwitchOff');

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
   * Staff have two more button to force switch
   */
   render() {
    return (
      <div className="machine used">
        <div className="row">
  
          <div className="col-xs-6">
  
            <div className="machine-action-info">
              <div className="machine-info-content">
                <div className="machine-time-label">
                  Usage time
                </div>
                <Timer activation={this.props.activation}/>
              </div>
            </div>
  
          </div>
  
          <div className="col-xs-6">
  
          { this.props.isStaff ? (
            
            <table className="machine-activation-table">
              <tr>
                <td rowSpan="2">
                  <button
                    className="btn btn-lg btn-danger btn-block"
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
              className="btn btn-lg btn-danger btn-block"
              onClick={this.endActivation}>
              Stop
            </button>
  
          )}
  
          </div>
  
        </div>
      </div>
    );
  }
});

export default BusyMachine;

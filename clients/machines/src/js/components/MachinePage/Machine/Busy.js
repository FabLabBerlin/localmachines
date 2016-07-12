var React = require('react');
var RepairButton = require('../../Feedback/RepairButton');
var Timer = require('./Timer');

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
   */
   render() {
    return (
      <div className="machine used">
        <div className="row">
  
          <div className="col-xs-6">
  
            <div className="machine-action-info">
              <div className="machine-info-content">
                <Timer activation={this.props.activation}/>
              </div>
            </div>
  
          </div>
  
          <div className="col-xs-6">
  
          { this.props.isStaff ? (
            
            <table className="machine-activation-table">
              <tbody>
                <tr>
                  <td rowSpan="2">
                    <button
                      className="btn btn-lg btn-danger btn-block"
                      onClick={this.endActivation}>
                      Stop
                    </button>
                  </td>
                </tr>
              </tbody>
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

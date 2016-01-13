var React = require('react');
var ForceSwitchOn = require('./ForceSwitchOn');
var ForceSwitchOff = require('./ForceSwitchOff');
var MachineActions = require('../../../modules/Machine/actions');

/*
 * Div displayed the machine is free
 * Can activate an activation
 */
var FreeMachine = React.createClass({

  /*
   * Try to activate the machine
   */
  startActivation() {
    this.props.func();
  },

  /*
   * Render Free machine div
   * If the machine is free, the component will be displayed
   * If is admin, two button will also be displayed
   */
  render() {
    var imageUrl;
    if (this.props.machine && this.props.machine.Image) {
      imageUrl = '/files/' + this.props.machine.Image;
    } else {
      imageUrl = '/machines/img/img-machine-placeholder.svg';
    }

    return (
      <div className="machine available">
        <div className="row">
          <div className="col-xs-6">
  
            {this.props.activation}
            <div className="machine-action-info">
              <img className="machine-image" src={imageUrl}/>
            </div>
          
          </div>
          <div className="col-xs-6">
  
            { this.props.isAdmin ? (
  
              <table className="machine-activation-table">
                <tr>
                  <td rowSpan="2">
                    <button
                      className="btn btn-lg btn-primary btn-block"
                      onClick={this.startActivation}>
                      Start
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
  
            ) : (
  
              <button
                className="btn btn-lg btn-primary btn-block"
                onClick={this.startActivation}>
                Start
              </button>
            
            )}
            
          </div>
        </div>
      </div>
    );
  }
});

export default FreeMachine;

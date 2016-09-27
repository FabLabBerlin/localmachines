var React = require('react');
var MachineActions = require('../../../actions/MachineActions');

/*
 * Div displayed the machine is free
 * Can activate an activation
 */
var FreeMachine = React.createClass({

  startActivation() {
    this.props.func();
  },

  render() {
    var imageUrl;
    if (this.props.machine && this.props.machine.ImageSmall) {
      imageUrl = '/files/' + this.props.machine.ImageSmall;
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
  
            { this.props.isStaff ? (
  
              <table className="machine-activation-table">
                <tbody>
                  <tr>
                    <td rowSpan="2">
                      <button
                        className="btn btn-lg btn-primary btn-block"
                        onClick={this.startActivation}>
                        Start
                      </button>
                    </td>
                  </tr>
                </tbody>
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

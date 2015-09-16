var React = require('react');
var ForceSwitchOn = require('./ForceSwitchOn');
var ForceSwitchOff = require('./ForceSwitchOff');
var MachineActions = require('../actions/MachineActions');
var MaintenanceSwitch = require('./MaintenanceSwitch');
var RepairButton = require('./Feedback/RepairButton');

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
    if (this.props.info && this.props.info.Image) {
      imageUrl = '/files/' + this.props.info.Image;
    } else {
      imageUrl = '/machines/img/img-machine-placeholder.svg';
    }

    return (
      <div>
        <div className="row">
          <div className="col-xs-6">
            {this.props.activation}
            <div className="machine-action-info">
              <img className="machine-image"
                   src={imageUrl}/>
            </div>
          </div>
          <div className="col-xs-6">
            <button
              className="btn btn-lg btn-primary btn-block"
              onClick={this.startActivation}
              >Start </button>
          </div>
        </div>
        <ul className="machine-extra-actions">
          
          {this.props.isAdmin ? (
            <li className="action-item">
              <ForceSwitchOn force={this.props.force}/>
            </li>
          ) : ''}

          {this.props.isAdmin ? (
            <li className="action-item">
              <ForceSwitchOff force={this.props.force}/>
            </li>
          ) : ''}

          <li className="action-item">
            <MaintenanceSwitch machineId={this.props.info.Id}/>
          </li>
          <li className="action-item">
            <RepairButton machineId={this.props.info.Id}/>
          </li>
        </ul>
      </div>
    );
  }
});

export default FreeMachine;

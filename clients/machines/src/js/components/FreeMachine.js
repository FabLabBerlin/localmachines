var React = require('react');
var ForceSwitch = require('./ForceSwitch');
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
              <ForceSwitch isAdmin={this.props.isAdmin} force={this.props.force}/>
          </div>
        </div>
        <ul className="machine-extra-actions">
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

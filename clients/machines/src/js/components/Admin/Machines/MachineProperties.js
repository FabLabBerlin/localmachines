import MachineActions from '../../../actions/MachineActions';
import React from 'react';


var MachineProperties = React.createClass({
  render() {
    const machine = this.props.machine;
    const machineType = {};

    return (
      <div>
        <div className="row">
          <div className="col-sm-3" ng-show="data.location.FeatureSpaces">
            <div className="form-group">
              <label>Start time (seconds)</label>
              <input type="number"
                     onChange={this.update.bind(this, 'GracePeriod')}
                     className="form-control"
                     value={machine.get('GracePeriod')} />
            </div>
          </div>

          <div className="col-sm-3">
            <div className="form-group">
              <label>
                <input type="checkbox"
                       onClick={this.update.bind(this, 'Visible')}
                       checked={machine.get('Visible')} /> Visible
              </label>
            </div>
          </div>
        </div>

        <div className="row">
          <div className="col-sm-3">
            <div className="form-group">
              <label>Machine Type</label>
              <select className="form-control"
                      onChange={this.update.bind(this, 'TypeId')}
                      value={machine.get('TypeId')}>
                <option value="0" selected disabled>Select type</option>
                <option ng-repeat="machineType in machineTypes"
                        value="{machineType.Id}">
                  {machineType.Name}
                </option>
              </select>
            </div>
          </div>

          <div className="col-sm-3">
            <div className="form-group">
              <label>Machine Brand</label>
              <input type="text"
                     className="form-control"
                     onChange={this.update.bind(this, 'Brand')}
                     placeholder="Enter machine brand"
                     value={machine.get('Brand')} />
            </div>
          </div>
        </div>

        <div className="row">
          <div className="col-sm-6">
            <div className="form-group">
              <label>Dimensions</label>
              <input type="text"
                     className="form-control"
                     onChange={this.update.bind(this, 'Dimensions')}
                     placeholder="Enter dimensions"
                     value={machine.get('Dimensions')} />
            </div>
          </div>

          <div className="col-sm-6">
            <div className="form-group">
              <label>Workspace Dimensions</label>
              <input type="text"
                     className="form-control"
                     onChange={this.update.bind(this, 'WorkspaceDimensions')}
                     placeholder="E.g. 200 mm x 200 mm x 200 mm or 1.5 m x 3 m"
                     value={machine.get('WorkspaceDimensions')} />
            </div>
          </div>
        </div>
      </div>
    );
  },

  update(name, e) {
    const id = this.props.machine.get('Id');
    var value = e.target.value;
    switch (name) {
    case 'GracePeriod':
      value = parseInt(value);
      break;
    case 'Visible':
      value = e.target.checked;
      break;
    }
    MachineActions.updateMachineField(id, name, value);
  }
});

export default MachineProperties;

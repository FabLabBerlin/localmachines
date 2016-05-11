var _ = require('lodash');
var $ = require('jquery');
var getters = require('../../../getters');
var LoaderLocal = require('../../LoaderLocal');
var LocationGetters = require('../../../modules/Location/getters');
var Machines = require('../../../modules/Machines');
var MachineActions = require('../../../actions/MachineActions');
var Navigation = require('react-router').Navigation;
var React = require('react');
var reactor = require('../../../reactor');
var toastr = require('../../../toastr');


var FirstRow = React.createClass({
  render() {
    const machine = this.props.machine;

    return (
      <div className="row">

        <div className="col-sm-3">
          <div className="form-group">
            <label>Machine Name</label>
            <input id="machine-name"
                   type="text"
                   className="form-control"
                   onChange={this.update.bind(this, 'Name')}
                   placeholder="Enter machine name"
                   value={machine.get('Name')}/>
          </div>
        </div>

        <div className="col-sm-3">
          <div className="form-group">
            <label>Short Name</label>
            <input type="text"
                   onChange={this.update.bind(this, 'Shortname')}
                   className="form-control"
                   placeholder="Enter short name"
                   value={machine.get('Shortname')}/>
          </div>
        </div>

        <div className="col-sm-3">
          <div className="form-group">
            <label>Price</label>
            <input type="text"
                   onChange={this.update.bind(this, 'Price')}
                   className="form-control"
                   placeholder="Enter price"
                   value={machine.get('Price')}/>
          </div>
        </div>

        <div className="col-sm-3">
          <div className="form-group">
            <label>Price Unit</label>
            <select className="form-control"
                    onChange={this.update.bind(this, 'PriceUnit')}
                    id="machine-price-unit"
                    value={machine.get('PriceUnit')}>
              <option value="" disabled>Select unit</option>
              <option value="minute">minute</option>
              <option value="hour">hour</option>
              <option value="day">day</option>
            </select>
          </div>
        </div>

      </div>
    );
  },

  update(name, e) {
    const id = this.props.machine.get('Id');
    var value = e.target.value;
    switch (name) {
    case 'Price':
      value = parseFloat(value);
      break;
    }
    MachineActions.updateMachineField(id, name, value);
  }
});


var SecondRow = React.createClass({
  render() {
    const machine = this.props.machine;

    var machineImageFile;
    var machineImageNewFile;
    var machineImageNewFileName;
    var machineImageNewFileSize;

    if (machine.get('Image')) {
      machineImageFile = '/files/' + machine.get('Image');
    }

    return (
      <div className="row">

        <div className="col-sm-3">
          <div className="form-group">
            <label>Machine Description</label>
            <textarea className="form-control"
                      onChange={this.update.bind(this, 'Description')}
                      placeholder="Enter machine description"
                      value={machine.get('Description')}
                      rows="5" />
          </div>
        </div>

        <div className="col-sm-6">
          <label>Image</label>
          <div className="row"> 
            <div className="col-sm-6">
              <div className="form-group">
                <img id="machine-image"
                  src={machineImageNewFile || machineImageFile || 'assets/img/img-machine-placeholder.svg'}
                  alt="Machine image"/>
              </div>
            </div>
            <div className="col-sm-6">
              <div className="form-group">
                <input type="file"
                       onchange="angular.element(this).scope().machineImageLoad(this)"/>
                <button className="btn btn-primary btn-block"
                        ng-disabled="!machineImageNewFile"
                        ng-click="machineImageReplace()">
                  <i className="fa fa-file-image-o"/>&nbsp;Replace
                </button>
              </div>
              <p>Supported types: svg</p>
              <p>Name: {machineImageNewFileName || machine.get('ImageName')}<br/>Size: {machineImageNewFileSize || machine.get('ImageSize')}</p>
            </div>
          </div>

        </div>

      </div>
    );
  },

  update(name, e) {
    const id = this.props.machine.get('Id');
    console.log('updating', name, 'with', e.target.value);
    MachineActions.updateMachineField(id, name, e.target.value);
  }
});


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


var NetswitchConfig = React.createClass({
  applyConfig() {
    toastr.error('Not implemented here yet');
  },

  render() {
    const machine = this.props.machine;

    var netswitchConfigStatus;

    return (
      <div className="row">

        <div className="col-sm-12">
          <label>NetSwitch Config</label>
          <div className="row">

            {machine.get('NetswitchType') ?
              (
                <div className="col-sm-3">
                  <div className="form-group">
                    <input type="text"
                           className="form-control"
                           placeholder="Host (mfi only)"
                           value={machine.get('NetswitchHost')}
                           onChange={this.update.bind(this, 'NetswitchHost')}/>
                  </div>
                </div>
              ) : (
                <div>
                  <div className="col-sm-3">
                    <div className="form-group">
                      <input type="text"
                             className="form-control"
                             placeholder="On URL"
                             value={machine.get('NetswitchUrlOn')}
                             onChange={this.update.bind(this, 'NetswitchUrlOn')}/>
                    </div>
                  </div>
                  <div className="col-sm-3">
                    <div className="form-group">
                      <input type="text"
                             className="form-control"
                             placeholder="Off URL"
                             value={machine.get('NetswitchUrlOff')}
                             onChange={this.update.bind(this, 'NetswitchUrlOff')}/>
                    </div>
                  </div>
                </div>
              )
            }

          </div>
          <div className="row">

            <div className="col-sm-3">
              <div className="form-group">
                <select className="form-control"
                        value={machine.get('NetswitchType')}
                        onChange={this.update.bind(this, 'NetswitchType')}>
                  <option value="">Custom Powerswitch</option>
                  <option value="mfi">Ubiquiti mFi Powerswitch</option>
                </select>
              </div>
            </div>

            <div className="col-sm-3">
              <div className="form-group">
                {netswitchConfigStatus ?
                  (
                    <button className="btn btn-danger btn-block"
                            disabled="true"
                            type="button">
                      <i className="fa fa-refresh fa-spin"></i>
                      {netswitchConfigStatus}
                    </button>
                  ) : null
                }
                {(!netswitchConfigStatus && machine.NetswitchType === 'mfi') ?
                  (
                    <button className="btn btn-danger btn-block"
                            onClick={this.applyConfig}
                            type="button">
                      Upgrade Powerswitch
                    </button>
                  ) : null
                }
              </div>
            </div>
            
          </div>
        </div>

      </div>
    );
  },

  update(name, e) {
    const id = this.props.machine.get('Id');
    console.log('updating', name, 'with', e.target.value);
    MachineActions.updateMachineField(id, name, e.target.value);
  }
});


var Buttons = React.createClass({
  render() {
    const machine = this.props.machine;

    return (
      <div className="pull-right">

        {machine.get('Archived') ? (
          <button className="btn btn-danger"
                  onClick={this.toggleArchived}>
            <i className="fa fa-archive"></i>&nbsp;Unarchive
          </button>
        ) : (
          <button className="btn btn-danger"
                  onClick={this.toggleArchived}>
            <i className="fa fa-archive"></i>&nbsp;Archive
          </button>
        )}

        <button className="btn btn-primary"
                onClick={this.save}>
          <i className="fa fa-save"></i>&nbsp;Save
        </button>

      </div>
    );
  },

  save() {
    var machine = this.props.machine.toJS();

    console.log('machine:', machine);

    $.ajax({
      method: 'PUT',
      url: '/api/machines/' + machine.Id,
      contentType: 'application/json',
      data: JSON.stringify(machine),
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      toastr.success('Update successful');
    })
    .error(function(message, statusCode) {
      if (statusCode === 400 && message.indexOf('Dimensions') >= 0) {
        toastr.error(message);
      } else if (statusCode === 400 && message.indexOf('Found machine with same netswitch host') >= 0) {
        toastr.error(message);
      } else {
        toastr.error('Failed to update');
      }
    });
  },

  toggleArchived() {
    const machine = this.props.machine;
    const action = machine.get('Archived') ? 'unarchived' : 'archived';

    $.ajax({
      method: 'POST',
      url: '/api/machines/' + machine.get('Id') + '/set_archived?archived=' + !machine.get('Archived')
    })
    .success(function(data) {
      toastr.info('Successfully ' + action + ' machine');
      const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
      const uid = reactor.evaluateToJS(getters.getUid);
      MachineActions.apiGetUserMachines(locationId, uid);
    })
    .error(function() {
      toastr.error('Failed to ' + action + ' machine');
    });
  }
});


var Machine = React.createClass({

  mixins: [ Navigation, reactor.ReactMixin ],

  getDataBindings() {
    return {
      location: LocationGetters.getLocation,
      machines: Machines.getters.getMachines
    };
  },

  componentWillMount() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);
    MachineActions.apiGetUserMachines(locationId, uid);
  },

  render() {
    const machineId = parseInt(this.props.params.machineId);
    var machine;
    if (this.state.machines) {
      machine = this.state.machines.find((m) => {
        return m.get('Id') === machineId;
      });
    }
    if (machine) {
      return (
        <div className="container-fluid">
          <h1>Edit Machine</h1>

          <hr />

          <FirstRow machine={machine} />
          <SecondRow machine={machine} />
          <MachineProperties machine={machine} />
          <NetswitchConfig machine={machine} />

          <Buttons machine={machine} />
        </div>
      );
    } else {
      return <LoaderLocal/>;
    }
  }

});

export default Machine;

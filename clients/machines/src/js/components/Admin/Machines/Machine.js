var _ = require('lodash');
var getters = require('../../../getters');
var LoaderLocal = require('../../LoaderLocal');
var LocationGetters = require('../../../modules/Location/getters');
var MachineActions = require('../../../actions/MachineActions');
var Navigation = require('react-router').Navigation;
var React = require('react');
var reactor = require('../../../reactor');


var FirstRow = React.createClass({
  render() {
    const machine = this.props.machine;

    return (
      <div className="row">

        <div className="col-sm-3">
          <div className="form-group">
            <label>Machine Name</label>
            <input type="text"
                   className="form-control"
                   placeholder="Enter machine name"
                   defaultValue={machine.get('Name')}/>
          </div>
        </div>

        <div className="col-sm-3">
          <div className="form-group">
            <label>Short Name</label>
            <input type="text"
                   className="form-control"
                   placeholder="Enter short name"
                   defaultValue={machine.get('Shortname')}/>
          </div>
        </div>

        <div className="col-sm-3">
          <div className="form-group">
            <label>Price</label>
            <input type="text"
                   className="form-control"
                   placeholder="Enter price"
                   defaultValue={machine.get('Price')}/>
          </div>
        </div>

        <div className="col-sm-3">
          <div className="form-group">
            <label>Price Unit</label>
            <select className="form-control"
                    defaultValue={machine.get('PriceUnit')}>
              <option value="" disabled>Select unit</option>
              <option value="minute">minute</option>
              <option value="hour">hour</option>
              <option value="day">day</option>
            </select>
          </div>
        </div>

      </div>
    );
  }
});


var SecondRow = React.createClass({
  render() {
    const machine = this.props.machine;

    var machineImageFile;
    var machineImageNewFile;
    var machineImageNewFileName;
    var machineImageNewFileSize;

    if (machine.Image) {
      machineImageFile = '/files/' + machine.Image;
    }

    return (
      <div className="row">

        <div className="col-sm-3">
          <div className="form-group">
            <label>Machine Description</label>
            <textarea className="form-control"
                      placeholder="Enter machine description"
                      defaultValue={machine.get('Description')}
                      rows="5">
            </textarea>
          </div>
        </div>

        <div className="col-sm-6">
          <label>Image</label>
          <div className="row"> 
            <div className="col-sm-6">
              <div className="form-group">
                <img id="machine-image"
                  src="{machineImageNewFile || machineImageFile || 'assets/img/img-machine-placeholder.svg'}"
                  alt="Machine image"/>
              </div>
            </div>
            <div className="col-sm-6">
              <div className="form-group">
                <input type="file" onchange="angular.element(this).scope().machineImageLoad(this)"/>
                <button className="btn btn-primary btn-block" ng-disabled="!machineImageNewFile" ng-click="machineImageReplace()">
                  <i className="fa fa-file-image-o"></i>&nbsp;Replace
                </button>
              </div>
              <p>Supported types: svg</p>
              <p>Name: {machineImageNewFileName || machine.ImageName}<br/>Size: {machineImageNewFileSize || machine.ImageSize}</p>
            </div>
          </div>

        </div>

      </div>
    );
  }
});


var MachineProperties = React.createClass({
  render() {
    const machineType = {};

    return (
      <div>
        <div className="row">
          <div className="col-sm-3" ng-show="data.location.FeatureSpaces">
            <div className="form-group">
              <label>Start time (seconds)</label>
              <input type="number"
                     id="grace-period"
                     className="form-control"
                     ng-model="machine.GracePeriod" />
            </div>
          </div>

          <div className="col-sm-3">
            <div className="form-group">
              <label>
                <input type="checkbox" ng-model="machine.Visible" /> Visible
              </label>
            </div>
          </div>
        </div>

        <div className="row">
          <div className="col-sm-3">
            <div className="form-group">
              <label>Machine Type</label>
              <select className="form-control" id="machine-type" 
                      ng-model="machine.TypeId">
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
              <input type="text" id="machine-brand" className="form-control"
                     placeholder="Enter machine brand" ng-model="machine.Brand" />
            </div>
          </div>
        </div>

        <div className="row">
          <div className="col-sm-6">
            <div className="form-group">
              <label>Dimensions</label>
              <input type="text" id="machine-dimensions" className="form-control"
                     placeholder="Enter dimensions" ng-model="machine.Dimensions" />
            </div>
          </div>

          <div className="col-sm-6">
            <div className="form-group">
              <label>Workspace Dimensions</label>
              <input type="text" id="machine-workspace-dimensions" className="form-control"
                     placeholder="E.g. 200 mm x 200 mm x 200 mm or 1.5 m x 3 m" ng-model="machine.WorkspaceDimensions" />
            </div>
          </div>
        </div>
      </div>
    );
  }
});


var NetswitchConfig = React.createClass({
  render() {
    const machine = this.props.machine;

    var netswitchConfigStatus;

    return (
      <div className="row">

        <div className="col-sm-12">
          <label>NetSwitch Config</label>
          <div className="row">

            <div className="col-sm-3">
              <div className="form-group" ng-hide="machine.NetswitchType">
                <input type="text"
                       className="form-control"
                       placeholder="On URL"
                       ng-model="machine.NetswitchUrlOn"
                       ng-change="registerUnsavedChange()"/>
              </div>
              <div className="form-group" ng-show="machine.NetswitchType">
                <input type="text"
                       className="form-control"
                       placeholder="Host (mfi only)"
                       ng-model="machine.NetswitchHost"
                       ng-change="registerUnsavedChange()"/>
              </div>
            </div>

            <div className="col-sm-3">
              <div className="form-group" ng-hide="machine.NetswitchType">
                <input type="text"
                       className="form-control"
                       placeholder="Off URL"
                       ng-model="machine.NetswitchUrlOff"
                       ng-change="registerUnsavedChange()"/>
              </div>
            </div>

          </div>
          <div className="row">

            <div className="col-sm-3">
              <div className="form-group">
                <select className="form-control"
                        ng-model="machine.NetswitchType"
                        ng-change="registerUnsavedChange()">
                  <option value="" selected>Custom Powerswitch</option>
                  <option value="mfi">Ubiquiti mFi Powerswitch</option>
                </select>
              </div>
            </div>

            <div className="col-sm-3">
              <div className="form-group">
                {netswitchConfigStatus ?
                  (
                    <button className="btn btn-danger btn-block"
                            ng-click="applyConfig()"
                            ng-disabled="true"
                            ng-show="netswitchConfigStatus"
                            type="button">
                      <i className="fa fa-refresh fa-spin"></i>
                      {netswitchConfigStatus}
                    </button>
                  ) :
                  (
                    <button className="btn btn-danger btn-block"
                            ng-click="applyConfig()"
                            ng-show="machine.NetswitchType === 'mfi' && !netswitchConfigStatus"
                            type="button">
                      Upgrade Powerswitch
                    </button>
                  )
                }
              </div>
            </div>
            
          </div>
        </div>

      </div>
    );
  }
});


var Machine = React.createClass({

  mixins: [ Navigation, reactor.ReactMixin ],

  /*
   * If not logged then redirect to the login page
   */
  statics: {
    willTransitionTo(transition) {
      const isLogged = reactor.evaluateToJS(getters.getIsLogged);
      if(!isLogged) {
        transition.redirect('login');
      }
    }
  },

  getDataBindings() {
    return {
      location: LocationGetters.getLocation,
      machines: getters.getMachines
    };
  },

  componentWillMount() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);
    MachineActions.apiGetUserMachines(locationId, uid);
  },

  render() {
    const machineId = parseInt(this.props.params.machineId);
    var machines = this.state.machines.toJS();
    if (this.state.machines) {
      var machine = this.state.machines.find((m) => {
        return m.get('Id') === machineId;
      });

      return (
        <div className="container-fluid">
          <h1>Edit Machine</h1>

          <hr />

          <FirstRow machine={machine} />
          <SecondRow machine={machine} />
          <MachineProperties machine={machine} />
          <NetswitchConfig machine={machine} />
        </div>
      );
    } else {
      return <LoaderLocal/>;
    }
  }

});

export default Machine;

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
        </div>
      );
    } else {
      return <LoaderLocal/>;
    }
  }

});

export default Machine;

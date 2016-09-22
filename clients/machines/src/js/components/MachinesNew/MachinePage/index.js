var getters = require('../../../getters');
var LoaderLocal = require('../../LoaderLocal');
var LocationGetters = require('../../../modules/Location/getters');
var MachineActions = require('../../../actions/MachineActions');
var Machines = require('../../../modules/Machines');
var React = require('react');
var reactor = require('../../../reactor');
var toastr = require('../../../toastr');


var MachinePage = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      locationId: LocationGetters.getLocationId,
      machines: Machines.getters.getMachines
    };
  },

  componentWillMount() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);
    MachineActions.apiGetUserMachines(locationId, uid);
  },

  imgUrl() {
    const m = this.machine();

    if (m.get('Image')) {
      return '/files/' + m.get('Image');
    } else {
      return '/machines/img/img-machine-placeholder.svg';
    }
  },

  machine() {
    const machineId = parseInt(this.props.params.machineId);
    var m;

    if (this.state.machines) {
      m = this.state.machines.find((mm) => {
        return mm.get('Id') === machineId;
      });
    }

    return m;
  },

  render() {
    const m = this.machine();

    if (!m) {
      return <LoaderLocal/>;
    }

    return (
      <div className="container-fluid">
        <div id="m-header">
          <h2>{m.get('Name')} ({m.get('Brand')})</h2>
          <img src={this.imgUrl()}/>
          <div id="m-header-panel">
            <div id="m-start-stop">
              START
            </div>
          </div>
        </div>
      </div>
    );
  }

});

export default MachinePage;

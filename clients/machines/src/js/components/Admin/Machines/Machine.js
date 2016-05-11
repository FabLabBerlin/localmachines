var _ = require('lodash');
var BasicData = require('./BasicData');
var Buttons = require('./Buttons');
var getters = require('../../../getters');
var ImageUpload = require('./ImageUpload');
var LoaderLocal = require('../../LoaderLocal');
var LocationGetters = require('../../../modules/Location/getters');
var Machines = require('../../../modules/Machines');
var MachineActions = require('../../../actions/MachineActions');
var MachineProperties = require('./MachineProperties');
var Navigation = require('react-router').Navigation;
var NetswitchConfig = require('./NetswitchConfig');
var React = require('react');
var reactor = require('../../../reactor');
var toastr = require('../../../toastr');


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

          <BasicData machine={machine} />
          <ImageUpload machine={machine} />
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

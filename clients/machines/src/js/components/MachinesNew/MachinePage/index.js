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

  render() {
    const machineId = parseInt(this.props.params.machineId);
    var m;

    if (this.state.machines) {
      m = this.state.machines.find((mm) => {
        return mm.get('Id') === machineId;
      });
    }
    if (!m) {
      return <LoaderLocal/>;
    }

    return (
      <div className="container-fluid">
        Individual machine page for {m.get('Name')} ({m.get('Brand')}) here...
      </div>
    );
  }

});

export default MachinePage;

var getters = require('../../getters');
var LoaderLocal = require('../LoaderLocal');
var LocationGetters = require('../../modules/Location/getters');
var MachineActions = require('../../actions/MachineActions');
var Machine = require('./Machine');
var Machines = require('../../modules/Machines');
var React = require('react');
var ReservationActions = require('../../actions/ReservationActions');
var reactor = require('../../reactor');
var toastr = require('../../toastr');
var UserActions = require('../../actions/UserActions');


var Section = React.createClass({
  render() {
    const machines = this.props.machines.sortBy(m => m.get('Name'));

    return (
      <div>
        <div className="ms-section-title">
          {this.title()}
        </div>
        <div>
          {machines.map((m, i) => <Machine key={i} machine={m}/>)}
        </div>
      </div>
    );
  },

  title() {
    var tid = this.props.typeId;
    var t = {
      0: 'Other',
      1: '3D Printers',
      2: 'CNC Mill',
      3: 'Heatpress',
      4: 'Knitting Machine',
      5: 'Lasercutters',
      6: 'Vinylcutter'
    }[tid];

    return t || tid;
  }
});


var MachinesPage = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      activations: Machines.getters.getActivations,
      locationId: LocationGetters.getLocationId,
      machines: Machines.getters.getMachines
    };
  },

  componentWillMount() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);

    ReservationActions.load();
    UserActions.fetchUser(uid);

    if (window.WebSocket) {
      MachineActions.wsDashboard(this.context.router, locationId);
    } else {
      MachineActions.lpDashboard(this.context.router, locationId);
    }
  },

  render() {
    if (!this.state.locationId || !this.state.machines ||
        this.state.machines.toList().count() === 0 ||
        !this.state.activations) {
      return <LoaderLocal/>;
    }

    const activationsByMachine = this.state.activations
    .groupBy(a => a.get('MachineId'));

    const machinesByType = this.state.machines
    .toList()
    .map(m => {
      const as = activationsByMachine.get(m.get('Id'));
      return as ? m.set('activation', as.get(0)) : m;
    })
    .filter(m => {
      return m.get('LocationId') === this.state.locationId &&
        !m.get('Archived');
    })
    .groupBy(m => m.get('TypeId'));

    return (
      <div id="ms" className="container-fluid">
        <Section typeId={1} machines={machinesByType.get(1)}/>
        <Section typeId={2} machines={machinesByType.get(2)}/>
        <Section typeId={3} machines={machinesByType.get(3)}/>
        <Section typeId={4} machines={machinesByType.get(4)}/>
        <Section typeId={5} machines={machinesByType.get(5)}/>
        <Section typeId={6} machines={machinesByType.get(6)}/>
        <Section typeId={0} machines={machinesByType.get(0)}/>
      </div>
    );
  }

});

export default MachinesPage;

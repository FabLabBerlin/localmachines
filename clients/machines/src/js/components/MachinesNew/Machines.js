var _ = require('lodash');
var getters = require('../../getters');
var LoaderLocal = require('../LoaderLocal');
var LocationActions = require('../../actions/LocationActions');
var LocationGetters = require('../../modules/Location/getters');
var MachineActions = require('../../actions/MachineActions');
var Machine = require('./Machine');
var Machines = require('../../modules/Machines');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;
var React = require('react');
var ReservationActions = require('../../actions/ReservationActions');
var reactor = require('../../reactor');
var toastr = require('../../toastr');
var UserActions = require('../../actions/UserActions');


var Section = React.createClass({
  render() {
    if (!this.props.machines || this.props.machines.count() === 0) {
      return <div/>;
    }

    const machines = this.props.machines.sortBy(m => m.get('Name'));

    return (
      <div>
        <div className="ms-section-title">
          {this.props.title}
        </div>
        <div>
          {machines.map((m, i) => <Machine key={i} machine={m}/>)}
        </div>
      </div>
    );
  }
});


var MachinesPage = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      activations: Machines.getters.getActivations,
      locationId: LocationGetters.getLocationId,
      machines: Machines.getters.getMachines,
      machinesById: Machines.getters.getMachinesById,
      myMachines: Machines.getters.getMyMachines,
      upcomingReservations: getters.getUpcomingReservationsByMachineId
    };
  },

  componentWillMount() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);

    ReservationActions.load();
    UserActions.fetchUser(uid);
    LocationActions.loadUserLocations(uid);

    if (window.WebSocket) {
      MachineActions.wsDashboard(this.context.router, locationId);
    } else {
      MachineActions.lpDashboard(this.context.router, locationId);
    }
  },

  render() {
    const uid = reactor.evaluateToJS(getters.getUid);
    var myMachinesList;
    if (this.state.myMachines && this.state.upcomingReservations) {
      myMachinesList = [];
      const myMachineIds = this.state.myMachines.map(m => m.get('Id')).toJS();

      _.each(this.state.myMachines.toJS(), m => {
        myMachinesList.push(m);
      });

      _.each(this.state.upcomingReservations.toList().toJS(), r => {
        if (!_.includes(myMachineIds, r.MachineId) && r.UserId === uid) {
          const m = this.state.machinesById.get(r.MachineId);
          myMachinesList.push(m);
        }
      });
      myMachinesList = toImmutable(myMachinesList);
    }
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
        m.get('Visible') && !m.get('Archived');
    })
    .groupBy(m => m.get('TypeId'));

    if (myMachinesList) {
      myMachinesList = myMachinesList
        .map(m => {
          const as = activationsByMachine.get(m.get('Id'));
          return as ? m.set('activation', as.get(0)) : m;
        });
    }

    return (
      <div id="ms" className="container-fluid">
        {myMachinesList
          ? <Section title="My Machines" machines={myMachinesList}/>
          : null}
        <Section title="3D Printers" machines={machinesByType.get(1)}/>
        <Section title="CNC Mill" machines={machinesByType.get(2)}/>
        <Section title="Heatpress" machines={machinesByType.get(3)}/>
        <Section title="Knitting Machine" machines={machinesByType.get(4)}/>
        <Section title="Lasercutters" machines={machinesByType.get(5)}/>
        <Section title="Vinylcutter" machines={machinesByType.get(6)}/>
        <Section title="Other" machines={machinesByType.get(0)}/>
      </div>
    );
  }
});

export default MachinesPage;

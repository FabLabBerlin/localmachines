import _ from 'lodash';
import Categories from '../../modules/Categories';
import getters from '../../getters';
import LoaderLocal from '../LoaderLocal';
import Location from '../../modules/Location';
import MachineActions from '../../actions/MachineActions';
import Machine from './Machine';
import Machines from '../../modules/Machines';
import Nuclear from 'nuclear-js';
var toImmutable = Nuclear.toImmutable;
import React from 'react';
import ReservationActions from '../../actions/ReservationActions';
import reactor from '../../reactor';
import toastr from '../../toastr';
import UserActions from '../../actions/UserActions';


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
      categories: Categories.getters.getAll,
      locationId: Location.getters.getLocationId,
      machines: Machines.getters.getMachines,
      machinesById: Machines.getters.getMachinesById,
      myMachines: Machines.getters.getMyMachines,
      upcomingReservations: getters.getUpcomingReservationsByMachineId
    };
  },

  componentWillMount() {
    const locationId = reactor.evaluateToJS(Location.getters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);

    ReservationActions.load();
    UserActions.fetchUser(uid);
    Location.actions.loadUserLocations(uid);

    if (window.WebSocket) {
      MachineActions.wsDashboard(this.context.router, locationId);
    } else {
      MachineActions.lpDashboard(this.context.router, locationId);
    }

    Categories.actions.loadAll(locationId);
  },

  render() {
    const uid = reactor.evaluateToJS(getters.getUid);
    var myMachinesList;

    if (this.state.myMachines) {
      myMachinesList = [];
      const myMachineIds = this.state.myMachines.map(m => m.get('Id')).toJS();

      _.each(this.state.myMachines.toJS(), m => {
        myMachinesList.push(m);
      });

      if (this.state.upcomingReservations) {
        _.each(this.state.upcomingReservations.toList().toJS(), r => {
          if (!_.includes(myMachineIds, r.MachineId) && r.UserId === uid) {
            const m = this.state.machinesById.get(r.MachineId);
            myMachinesList.push(m);
          }
        });
      }

      myMachinesList = toImmutable(myMachinesList);
    }
    if (!this.state.locationId || !this.state.machines ||
        this.state.machines.toList().count() === 0 ||
        !this.state.activations ||
        !this.state.categories) {
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
        {this.state.categories.map(c => {
          return <Section key={c.get('Id')}
                          title={c.get('Name')}
                          machines={machinesByType.get(c.get('Id'))}/>;
        })}
      </div>
    );
  }
});

export default MachinesPage;

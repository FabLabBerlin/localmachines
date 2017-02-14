import Button from './Button';
import constants from '../constants';
import FeedbackDialogs from '../../Feedback/FeedbackDialogs';
import getters from '../../../getters';
import LoaderLocal from '../../LoaderLocal';
import Location from '../../../modules/Location';
import MachineActions from '../../../actions/MachineActions';
import MachineMixin from '../MachineMixin';
import Machines from '../../../modules/Machines';
import MaintenanceSwitch from './MaintenanceSwitch';
import OccupiedBy from './OccupiedBy';
import React from 'react';
import ReservationActions from '../../../actions/ReservationActions';
import ReservedBy from './ReservedBy';
import reactor from '../../../reactor';
import UpcomingReservation from './UpcomingReservation';
import UserActions from '../../../actions/UserActions';


var MachinePage = React.createClass({

  mixins: [ MachineMixin, reactor.ReactMixin ],

  getDataBindings() {
    return {
      activations: Machines.getters.getActivations,
      isStaff: Location.getters.getIsStaff,
      locationId: Location.getters.getLocationId,
      machines: Machines.getters.getMachines,
      reservationsByMachineId: getters.getActiveReservationsByMachineId,
      upcomingReservationsByMachineId: getters.getUpcomingReservationsByMachineId,
      user: getters.getUser,
      width: getters.getWidth
    };
  },

  componentWillMount() {
    const locationId = reactor.evaluateToJS(Location.getters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);
    MachineActions.apiGetUserMachines(locationId, uid);
    UserActions.fetchUser(uid);
    Location.actions.loadUserLocations(uid);
    ReservationActions.load();
    MachineActions.wsDashboard(null, locationId);
  },

  machine() {
    const machineId = parseInt(this.props.params.machineId);
    var m;

    if (this.state.machines) {
      m = this.state.machines.find((mm) => {
        return mm.get('Id') === machineId;
      });

      if (this.state.activations) {
        const as = this.state.activations
        .groupBy(a => a.get('MachineId'))
        .get(m.get('Id'));

        if (as) {
          m = m.set('activation', as.get(0));
        }
      }
    }

    return m;
  },

  render() {
    const m = this.machine();

    if (!m) {
      return <LoaderLocal/>;
    }

    var className;

    switch (this.status(true)) {
      case constants.AVAILABLE:
        className = 'm-start';
        break;
      case constants.LOCKED:
        className = 'm-locked';
        break;
      case constants.MAINTENANCE:
        className = 'm-maintenance';
        break;
      case constants.OCCUPIED:
        className = 'm-occupied';
        if (this.state.isStaff) {
          className += ' m-stop';
        }
        break;
      case constants.RESERVED:
        className = 'm-reserved';
        break;
      case constants.RUNNING:
        className = 'm-stop';
        break;
    }

    const small = this.state.width < 500;
    const style = {
      backgroundImage: 'url(' + this.imgUrl(small) + ')'
    };

    return (
      <div className="container-fluid">
        <div id="m-header">
          <h2>{m.get('Name')} ({m.get('Brand')})</h2>
          <div id="m-img" style={style}/>
          <div className={'row m-header-panel ' + className}>
            <div className="col-sm-4">
              <UpcomingReservation upcomingReservation={this.upcomingReservation()}/>
              <OccupiedBy activation={this.machine().get('activation')}/>
              <ReservedBy reservation={this.reservation()}/>
            </div>
            <div className="col-sm-4">
              <Button isStaff={this.state.isStaff}
                      machine={this.machine()}
                      reservation={this.reservation()}
                      status={this.status(true)}/>
              <br/>
              <MaintenanceSwitch.Off machine={m}/>
            </div>
            <div className="col-sm-4"/>
          </div>
          {this.status() !== constants.MAINTENANCE ?
            <div id="m-report" className="m-maintenance">
              <div className="m-maintenance-action" onClick={this.repair}>
                Report a machine failure
              </div>
              <br/>
              <MaintenanceSwitch.On machine={m}/>
            </div> : null
          }
        </div>
      </div>
    );
  },

  repair() {
    FeedbackDialogs.machineIssue(this.machine().get('Id'));
  }

});

export default MachinePage;

import Calendar from './Calendar';
import getters from '../../../../getters';
import Location from '../../../../modules/Location';
import MachineActions from '../../../../actions/MachineActions';
import NewReservation from '../../../Reservations/NewReservation';
import React from 'react';
import reactor from '../../../../reactor';
import ReservationActions from '../../../../actions/ReservationActions';
import ReservationRulesActions from '../../../../actions/ReservationRulesActions';
import Settings from '../../../../modules/Settings';
import UserActions from '../../../../actions/UserActions';


var ReservationPage = React.createClass({
  mixins: [ reactor.ReactMixin ],

  componentWillUnmount() {
    ReservationActions.newReservation.done();
  },

  clickCreate() {
    const mid = parseInt(this.props.params.machineId);
    ReservationActions.newReservation.create();
    ReservationActions.newReservation.setMachine({ mid });
    ReservationActions.newReservation.nextStep();
  },

  getDataBindings() {
    return {
      newReservation: getters.getNewReservation
    };
  },

  componentWillMount() {
    const locationId = reactor.evaluateToJS(Location.getters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);
    MachineActions.apiGetUserMachines(locationId, uid);
    UserActions.fetchUser(uid);
    Location.actions.loadUserLocations(uid);
    ReservationActions.load();
    ReservationRulesActions.load(locationId);
    MachineActions.wsDashboard(null, locationId);
    Settings.actions.loadSettings({locationId});
  },

  render() {
    const machineId = parseInt(this.props.params.machineId);

    if (this.state.newReservation) {
      return <NewReservation/>;
    } else {
      return (
        <div>
          <Calendar machineId={machineId}
                    clickCreate={this.clickCreate}/>
        </div>
      );
    }
  }
});

export default ReservationPage;

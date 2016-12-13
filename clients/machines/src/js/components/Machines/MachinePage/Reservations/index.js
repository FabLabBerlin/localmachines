var Calendar = require('./Calendar');
var getters = require('../../../../getters');
var Location = require('../../../../modules/Location');
var MachineActions = require('../../../../actions/MachineActions');
var NewReservation = require('../../../Reservations/NewReservation');
var React = require('react');
var reactor = require('../../../../reactor');
var ReservationActions = require('../../../../actions/ReservationActions');
var ReservationRulesActions = require('../../../../actions/ReservationRulesActions');
var Settings = require('../../../../modules/Settings');
var UserActions = require('../../../../actions/UserActions');


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

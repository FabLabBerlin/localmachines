var Calendar = require('./Calendar');
var getters = require('../../../../getters');
var LocationActions = require('../../../../actions/LocationActions');
var LocationGetters = require('../../../../modules/Location/getters');
var MachineActions = require('../../../../actions/MachineActions');
var React = require('react');
var reactor = require('../../../../reactor');
var ReservationActions = require('../../../../actions/ReservationActions');
var ReservationRulesActions = require('../../../../actions/ReservationRulesActions');
var UserActions = require('../../../../actions/UserActions');


var ReservationPage = React.createClass({
  componentWillMount() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);
    MachineActions.apiGetUserMachines(locationId, uid);
    UserActions.fetchUser(uid);
    LocationActions.loadUserLocations(uid);
    ReservationActions.load();
    ReservationRulesActions.load(locationId);
    MachineActions.wsDashboard(null, locationId);
  },

  render() {
    const machineId = parseInt(this.props.params.machineId);

    return (
      <div>
        <Calendar machineId={machineId}/>
      </div>
    );
  }
});

export default ReservationPage;

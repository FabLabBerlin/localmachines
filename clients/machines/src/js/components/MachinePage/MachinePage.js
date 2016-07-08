var _ = require('lodash');
var $ = require('jquery');
var getters = require('../../getters');
var LoaderLocal = require('../LoaderLocal');
var LocationGetters = require('../../modules/Location/getters');
var LocationActions = require('../../actions/LocationActions');
var LoginStore = require('../../stores/LoginStore');
var MachineList = require('./MachineList');
var MachineActions = require('../../actions/MachineActions');
var Machines = require('../../modules/Machines');
var NfcLogoutMixin = require('../Login/NfcLogoutMixin');
var LoginActions = require('../../actions/LoginActions');
var Navigation = require('react-router').Navigation;
var React = require('react');
var reactor = require('../../reactor');
var ReservationRulesActions = require('../../actions/ReservationRulesActions');
var ReservationActions = require('../../actions/ReservationActions');
var ScrollNav = require('../ScrollNav');
var toastr = require('../../toastr');
var UserActions = require('../../actions/UserActions');
var TutoringList = require('./TutoringList');

var MachinePage = React.createClass({

  mixins: [ Navigation, reactor.ReactMixin, NfcLogoutMixin ],

  getDataBindings() {
    return {
      user: getters.getUser,
      machines: Machines.getters.getMachines,
      activations: Machines.getters.getActivations,
      locations: LocationGetters.getLocations,
      location: LocationGetters.getLocation
    };
  },

  componentWillMount() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);
    UserActions.fetchUser(uid);
    ReservationActions.load();
    ReservationRulesActions.load(locationId);
    LocationActions.loadUserLocations(uid);
  },

  handleLogout() {
    LoginActions.logout(this.context.router);
  },

  componentDidMount() {
    LoginStore.onChangeLogout = this.onChangeLogout;
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    if (window.WebSocket) {
      MachineActions.wsDashboard(this.context.router, locationId);
    } else {
      MachineActions.lpDashboard(this.context.router, locationId);
    }
  },

  render() {
    console.log('MachinePage#render');
    if (this.state.activations && this.state.location && this.state.machines) {
      const locationTitle = this.state.location.get('Title');

      return (
        <div>
          <TutoringList />
          <MachineList
            user={this.state.user}
            machines={this.state.machines}
            activation={this.state.activations}
          />
          <div className="container-fluid">
            <button
              onClick={this.handleLogout}
              className="btn btn-lg btn-block btn-danger btn-logout-bottom">
              Sign out
            </button>
          </div>
          <ScrollNav/>
        </div>
      );
    } else {
      return <LoaderLocal/>;
    }
  }
});

export default MachinePage;

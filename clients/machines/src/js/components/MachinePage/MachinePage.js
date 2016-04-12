var _ = require('lodash');
var $ = require('jquery');
var getters = require('../../getters');
var LocationGetters = require('../../modules/Location/getters');
var LocationActions = require('../../actions/LocationActions');
var LoginStore = require('../../stores/LoginStore');
var MachineList = require('./MachineList');
var MachineStore = require('../../stores/MachineStore');
var MachineActions = require('../../actions/MachineActions');
var NfcLogoutMixin = require('../Login/NfcLogoutMixin');
var LoginActions = require('../../actions/LoginActions');
var Navigation = require('react-router').Navigation;
var React = require('react');
var reactor = require('../../reactor');
var ReservationRulesActions = require('../../actions/ReservationRulesActions');
var ReservationsActions = require('../../actions/ReservationsActions');
var ScrollNav = require('../ScrollNav');
var Location = require('./Location');
var toastr = require('../../toastr');
var UserActions = require('../../actions/UserActions');
var TutoringList = require('./TutoringList');

var MachinePage = React.createClass({

  mixins: [ Navigation, reactor.ReactMixin, NfcLogoutMixin ],

  /*
   * If not logged then redirect to the login page
   */
  statics: {
    willTransitionTo(transition) {
      const isLogged = reactor.evaluateToJS(getters.getIsLogged);
      if(!isLogged) {
        transition.redirect('login');
      }
    }
  },

  getDataBindings() {
    return {
      user: getters.getUser,
      machines: getters.getMachines,
      activations: getters.getActivations,
      locations: LocationGetters.getLocations,
      location: LocationGetters.getLocation
    };
  },

  /*
   * Start fetching the data
   * before the component is mounted
   */
  componentWillMount() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);
    UserActions.fetchUser(uid);
    MachineActions.apiGetUserMachines(locationId, uid);
    ReservationsActions.load();
    ReservationRulesActions.load(locationId);
    LocationActions.loadUserLocations(uid);
  },

  /*
   * Clear state while logout
   */
  clearState() {
    MachineActions.clearState();
  },

  /*
   * Logout with the exit button
   */
  handleLogout() {
    LoginActions.logout(this.context.router);
  },


  /*
   * Destructor
   * Stop the polling
   */
  componentWillUnmount() {
    this.nfcOnWillUnmount();
    this.clearState();
    clearInterval(this.interval);
  },

  /*
   * Call when the component is mounted in DOM
   * Synchronize invent = require(stores
   * Activate a polling (1,5s)
   */
  componentDidMount() {
    this.nfcOnDidMount();
    MachineStore.onChangeActivation = this.onChangeActivation;
    LoginStore.onChangeLogout = this.onChangeLogout;
    MachineStore.onChangeLogin = this.onChangeLogin;
    this.interval = setInterval(this.update, 1500);
  },

  /*
   * Render the user name
   * MachinList
   * exit button
   */
  render() {
    var machines = reactor.evaluateToJS(getters.getMachines);
    machines = _.sortBy(machines, (m) => {
      return m.Name;
    });
    if (this.state.activations && this.state.location) {
      var locationTitle = this.state.location.Title;
      return (
        <div>
          <div className="logged-user-name">
            <div className="text-center">
              <strong>Welcome to {locationTitle}</strong>
            </div>
          </div>
          <TutoringList />
          <MachineList
            user={this.state.user}
            machines={machines}
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
      return <div/>;
    }
  },

  /*
   * update
   *
   * Need polling for activation status and maintenance status
   */
  update() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    MachineActions.pollDashboard(this.context.router, locationId);
  }
});

export default MachinePage;

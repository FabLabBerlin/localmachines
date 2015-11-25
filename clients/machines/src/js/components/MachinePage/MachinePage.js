var $ = require('jquery');
var getters = require('../../getters');
var MachineList = require('./MachineList');
var LoginStore = require('../../stores/LoginStore');
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
      userInfo: getters.getUserInfo,
      machineInfo: getters.getMachineInfo,
      activationInfo: getters.getActivationInfo
    };
  },

  /*
   * Start fetching the data
   * before the component is mounted
   */
  componentWillMount() {
    const uid = reactor.evaluateToJS(getters.getUid);
    UserActions.getUserInfoFromServer(uid);
    MachineActions.apiGetUserMachines(uid);
    ReservationsActions.load();
    ReservationRulesActions.load();
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
    var machineInfo = reactor.evaluateToJS(getters.getMachineInfo);
    if (this.state.activationInfo) {
      return (
        <div>
          <div className="logged-user-name">
            <div className="text-center ng-binding">
              <i className="fa fa-user-secret"></i>&nbsp;
              {this.state.userInfo.get('FirstName')} {this.state.userInfo.get('LastName')}
            </div>
          </div>
          <TutoringList />
          <MachineList
            user={this.state.userInfo}
            info={machineInfo}
            activation={this.state.activationInfo}
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
    MachineActions.pollActivations();
    MachineActions.pollMachines();
  }
});

export default MachinePage;

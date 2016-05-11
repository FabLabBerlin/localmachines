var BillTable = require('./BillTable');
var getters = require('../../getters');
var LocationGetters = require('../../modules/Location/getters');
var LoginActions = require('../../actions/LoginActions');
var MachineActions = require('../../actions/MachineActions');
var Machines = require('../../modules/Machines');
var Membership = require('./Membership');
var NfcLogoutMixin = require('../Login/NfcLogoutMixin');
var {Navigation} = require('react-router');
var React = require('react');
var reactor = require('../../reactor');
var ScrollNav = require('../ScrollNav');
var SettingsActions = require('../../modules/Settings/actions');
var UserActions = require('../../actions/UserActions');


var SpendingsPage = React.createClass({

  /*
   * to use transitionTo/replaceWith/redirect and some function related to the router
   */
  mixins: [ Navigation, reactor.ReactMixin, NfcLogoutMixin ],

  /*
   * Fetching the user state from the store
   */
  getDataBindings() {
    return {
      user: getters.getUser,
      machines: Machines.getters.getMachines,
      bill: getters.getBill,
      memberships: getters.getMemberships
    };
  },

  componentDidMount() {
    this.nfcOnDidMount();
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);
    MachineActions.apiGetUserMachines(locationId, uid);
    UserActions.fetchUser(uid);
    UserActions.fetchBill(locationId, uid);
    UserActions.fetchMemberships(locationId, uid);
    SettingsActions.loadSettings({locationId});
  },

  componentWillUnmount() {
    this.nfcOnWillUnmount();
  },

  /*
   * Logout with the exit button
   */
  handleLogout() {
    LoginActions.logout();
  },

  render() {
    return (
      <div className="container">
        <h3>Your Memberships</h3>
        {<Membership memberships={this.state.memberships} />}

        <h3>Pay-As-You-Go</h3>
        <BillTable bill={this.state.bill} membership={this.state.membership}/>
        <ScrollNav/>
      </div>
    );
  }
});

export default SpendingsPage;

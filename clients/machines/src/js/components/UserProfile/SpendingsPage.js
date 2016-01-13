var BillTable = require('./BillTable');
var Login = require('../../modules/Login');
var Machine = require('../../modules/Machine');
var Membership = require('./Membership');
var NfcLogoutMixin = require('../Login/NfcLogoutMixin');
var {Navigation} = require('react-router');
var React = require('react');
var reactor = require('../../reactor');
var ScrollNav = require('../ScrollNav');
var User = require('../../modules/User');


var SpendingsPage = React.createClass({

  /*
   * to use transitionTo/replaceWith/redirect and some function related to the router
   */
  mixins: [ Navigation, reactor.ReactMixin, NfcLogoutMixin ],

  /*
   * If not logged then redirect to the login page
   */
  statics: {
    willTransitionTo(transition) {
      const isLogged = reactor.evaluateToJS(Login.getters.getIsLogged);
      if(!isLogged) {
        transition.redirect('login');
      }
    }
  },

  /*
   * Fetching the user state from the store
   */
  getDataBindings() {
    return {
      user: User.getters.getUser,
      machines: Machine.getters.getMachines,
      bill: User.getters.getBill,
      memberships: User.getters.getMemberships
    };
  },

  componentDidMount() {
    this.nfcOnDidMount();
    const uid = reactor.evaluateToJS(getters.getUid);
    Machine.actions.apiGetUserMachines(uid);
    User.actions.fetchUser(uid);
    User.actions.fetchBill(uid);
    User.actions.fetchMemberships(uid);
  },

  componentWillUnmount() {
    this.nfcOnWillUnmount();
  },

  /*
   * Logout with the exit button
   */
  handleLogout() {
    Login.actions.logout();
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

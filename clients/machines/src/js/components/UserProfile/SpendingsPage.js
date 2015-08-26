var BillTable = require('./BillTable');
var getters = require('../../getters');
var MachineActions = require('../../actions/MachineActions');
var Membership = require('./Membership');
var {Navigation} = require('react-router');
var React = require('react');
var reactor = require('../../reactor');
var UserActions = require('../../actions/UserActions');


var SpendingsPage = React.createClass({

  /*
   * to use transitionTo/replaceWith/redirect and some function related to the router
   */
  mixins: [ Navigation, reactor.ReactMixin ],

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

  /*
   * Fetching the user state from the store
   */
  getDataBindings() {
    return {
      userInfo: getters.getUserInfo,
      machineInfo: getters.getMachineInfo,
      billInfo: getters.getBillInfo,
      membershipInfo: getters.getMembership
    };
  },

  componentDidMount() {
    const uid = reactor.evaluateToJS(getters.getUid);
    MachineActions.apiGetUserMachines(uid);
    UserActions.getUserInfoFromServer(uid);
    UserActions.getInfoBillFromServer(uid);
    UserActions.getMembershipFromServer(uid);
  },

  render() {
    return (
      <div className="container">
        <h3>Your Monthly Spendings</h3>
        {<BillTable info={this.state.billInfo} />}

        <h3>Your Memberships</h3>
        {<Membership info={this.state.membershipInfo} />}
      </div>
    );
  }
});

export default SpendingsPage;

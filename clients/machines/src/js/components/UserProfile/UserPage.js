var getters = require('../../getters');
var MachineActions = require('../../actions/MachineActions');
var MachineList = require('./MachineList');
var {Navigation} = require('react-router');
var LoginActions = require('../../actions/LoginActions');
var NfcLogoutMixin = require('../NfcLogoutMixin');
var React = require('react');
var reactor = require('../../reactor');
var UserActions = require('../../actions/UserActions');
var UserForm = require('./UserForm');

/*
 * UserPage component:
 * manage the interaction with user
 * @children:
 *  - UserForm
 *  - MachineList
 *  - Membership
 */
var UserPage = React.createClass({

  /*
   * to use transitionTo/replaceWith/redirect and some function related to the router
   */
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
    this.nfcOnDidMount();
    const uid = reactor.evaluateToJS(getters.getUid);
    MachineActions.apiGetUserMachines(uid);
    UserActions.getUserInfoFromServer(uid);
    UserActions.getInfoBillFromServer(uid);
    UserActions.getMembershipFromServer(uid);
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

  /*
   * Submit the user information to the store via the action
   */
  handleSubmit() {
    const uid = reactor.evaluateToJS(getters.getUid);
    var userInfo = reactor.evaluateToJS(getters.getUserInfo);
    UserActions.submitState(uid, userInfo);
  },

  /*
   * When a change happend in the form:
   * @event: the event which occured
   * change the state to be coherent with the input values
   */
  handleChangeForm(event) {
    /*// Create a temporary state to replace the old one
    var tmpState = this.state.infoUser;
    tmpState[event.target.id] = event.target.value;
    this.setState({
      infoUser: tmpState
    });*/
    var key = event.target.id;
    var value = event.target.value;
    UserActions.setUserInfoProperty({ key, value });
  },

  /*
   * Send an action to update the password
   * @password: your new password
   */
  updatePassword(password) {
    const uid = reactor.evaluateToJS(getters.getUid);
    UserActions.updatePassword(uid, password);
  },

  /*
   * Render:
   *  - UserForm: form to update the user information
   *  - MachineList: machines the user can access
   *  - Membership: membership the user subscribe
   */
  render() {
    return (
      <div className="container">
        <h3>Your information</h3>
          <UserForm info={this.state.userInfo}
            func={this.handleChangeForm}
            passwordFunc={this.updatePassword}
            submit={this.handleSubmit}
          />
      </div>
    );
  }
});

export default UserPage;

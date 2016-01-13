var LoginGetters = require('../../modules/Login/getters');
var UserGetters = require('../../modules/User/getters');
var MachineActions = require('../../modules/Machine/actions');
var {Navigation} = require('react-router');
var LoginActions = require('../../modules/Login/actions');
var NfcLogoutMixin = require('../Login/NfcLogoutMixin');
var React = require('react');
var reactor = require('../../reactor');
var UserActions = require('../../modules/User/actions');
var UserForm = require('./UserForm');


var UserPage = React.createClass({

  mixins: [ Navigation, reactor.ReactMixin, NfcLogoutMixin ],

  /*
   * If not logged then redirect to the login page
   */
  statics: {
    willTransitionTo(transition) {
      const isLogged = reactor.evaluateToJS(LoginGetters.getIsLogged);
      if(!isLogged) {
        transition.redirect('login');
      }
    }
  },

  getDataBindings() {
    return {
      user: UserGetters.getUser
    };
  },

  componentDidMount() {
    this.nfcOnDidMount();
    const uid = reactor.evaluateToJS(LoginGetters.getUid);
    MachineActions.apiGetUserMachines(uid);
    UserActions.fetchUser(uid);
    UserActions.fetchBill(uid);
    UserActions.fetchMemberships(uid);
  },

  componentWillUnmount() {
    this.nfcOnWillUnmount();
  },

  handleLogout() {
    LoginActions.logout();
  },

  handleSubmit() {
    const uid = reactor.evaluateToJS(LoginGetters.getUid);
    var user = reactor.evaluateToJS(UserGetters.getUser);
    UserActions.updateUser(uid, user);
  },

  handleChangeForm(event) {
    var key = event.target.id;
    var value = event.target.value;
    UserActions.setUserProperty({ key, value });
  },

  updatePassword(password) {
    const uid = reactor.evaluateToJS(LoginGetters.getUid);
    UserActions.updatePassword(uid, password);
  },

  render() {
    return (
      <div className="container">
        <h3>Your information</h3>
          <UserForm user={this.state.user}
            func={this.handleChangeForm}
            passwordFunc={this.updatePassword}
            submit={this.handleSubmit}
          />
      </div>
    );
  }
});

export default UserPage;

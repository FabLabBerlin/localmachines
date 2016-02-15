var getters = require('../../getters');
var MachineActions = require('../../actions/MachineActions');
var {Navigation} = require('react-router');
var LoginActions = require('../../actions/LoginActions');
var NfcLogoutMixin = require('../Login/NfcLogoutMixin');
var React = require('react');
var reactor = require('../../reactor');
var UserActions = require('../../actions/UserActions');
var UserForm = require('./UserForm');


var UserPage = React.createClass({

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
      user: getters.getUser
    };
  },

  componentDidMount() {
    this.nfcOnDidMount();
    const locationId = reactor.evaluateToJS(getters.getLocation).Id;
    const uid = reactor.evaluateToJS(getters.getUid);
    MachineActions.apiGetUserMachines(uid);
    UserActions.fetchUser(uid);
    UserActions.fetchBill(locationId, uid);
    UserActions.fetchMemberships(locationId, uid);
  },

  componentWillUnmount() {
    this.nfcOnWillUnmount();
  },

  handleLogout() {
    LoginActions.logout();
  },

  handleSubmit() {
    const uid = reactor.evaluateToJS(getters.getUid);
    var user = reactor.evaluateToJS(getters.getUser);
    UserActions.updateUser(uid, user);
  },

  handleChangeForm(event) {
    var key = event.target.id;
    var value = event.target.value;
    UserActions.setUserProperty({ key, value });
  },

  updatePassword(password) {
    const uid = reactor.evaluateToJS(getters.getUid);
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

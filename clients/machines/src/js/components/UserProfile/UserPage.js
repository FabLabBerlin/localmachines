var getters = require('../../getters');
var LocationActions = require('../../actions/LocationActions');
var LocationGetters = require('../../modules/Location/getters');
var MachineActions = require('../../actions/MachineActions');
var LoginActions = require('../../actions/LoginActions');
var React = require('react');
var reactor = require('../../reactor');
var UserActions = require('../../actions/UserActions');
var UserForm = require('./UserForm');


var UserPage = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      user: getters.getUser
    };
  },

  componentDidMount() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);
    MachineActions.apiGetUserMachines(locationId, uid);
    UserActions.fetchUser(uid);
    UserActions.fetchBill(locationId, uid);
    UserActions.fetchMemberships(locationId, uid);
    LocationActions.loadUserLocations(uid);
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
        <h3>Your Information</h3>
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

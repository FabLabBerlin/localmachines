import getters from '../../getters';
import Location from '../../modules/Location';
import MachineActions from '../../actions/MachineActions';
import LoginActions from '../../actions/LoginActions';
import React from 'react';
import reactor from '../../reactor';
import UserActions from '../../actions/UserActions';
import UserForm from './UserForm';


var UserPage = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      user: getters.getUser
    };
  },

  componentDidMount() {
    const locationId = reactor.evaluateToJS(Location.getters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);
    MachineActions.apiGetUserMachines(locationId, uid);
    UserActions.fetchUser(uid);
    UserActions.fetchBill(locationId, uid);
    UserActions.fetchMemberships(locationId, uid);
    Location.actions.loadUserLocations(uid);
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

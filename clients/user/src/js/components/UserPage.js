import React from 'react';
import {Navigation} from 'react-router';
import UserForm from './UserForm';
import MachineList from './MachineList';
import BillTable from './BillTable';
import Membership from './Membership';
import UserActions from '../actions/UserActions';
import UserStore from '../stores/UserStore'

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
  mixins: [ Navigation ],

  /*
   * If not logged then redirect to the login page
   */
  statics: {
    willTransitionTo(transition) {
      if(!UserStore.getIsLogged()) {
        transition.redirect('login');
      }
    }
  },

  /*
   * Fetching the user state from the store
   */
  getInitialState() {
    return {
      infoUser: UserStore.getInfoUser(),
      infoMachine: UserStore.getInfoMachine(),
      infoBill: UserStore.getInfoBill(),
      infoMembership: UserStore.getMembership()
    };
  },

  /*
   * Submit the user information to the store via the action
   */
  handleSubmit() {
    UserActions.submitState(this.state.infoUser);
  },

  /*
   * When a change happend in the form:
   * @event: the event which occured
   * change the state to be coherent with the input values
   */
  handleChangeForm(event) {
    // Create a temporary state to replace the old one
    var tmpState = this.state.infoUser;
    tmpState[event.target.id] = event.target.value;
    this.setState({
      infoUser: tmpState
    });
  },

  /*
   * TODO: commentary
   */
  updatePassword(password) {
    UserActions.updatePassword(password);
  },

  /*
   * When logout, redirect to the login page
   */
  onChangeLogout() {
    if( !UserStore.getIsLogged() ){
      this.replaceWith('login');
    }
  },

  /*
   * To synchronize the logout call with the logout event
   */
  componentDidMount() {
    UserStore.onChangeLogout = this.onChangeLogout;
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
          <UserForm info={this.state.infoUser} 
            func={this.handleChangeForm}
            passwordFunc={this.updatePassword}
            submit={this.handleSubmit}
          />
          
        <h3>Machines you can use</h3>
        <MachineList info={this.state.infoMachine} />
          
        <h3>Your spending</h3>
        <BillTable info={this.state.infoBill} />

        <h3>Your Memberships</h3>
        <Membership info={this.state.infoMembership} />
      </div>
    );
  }
});

module.exports = UserPage;

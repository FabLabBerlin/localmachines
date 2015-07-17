import React from 'react'
import MachineList from './MachineList';
import MachineStore from '../stores/MachineStore';
import LoginActions from '../actions/LoginActions';
import {Navigation} from 'react-router';

var MachinePage = React.createClass({

  mixins: [ Navigation ],

  /*
   * If not logged then redirect to the login page
   */
  statics: {
    willTransitionTo(transition) {
      if(!MachineStore.getIsLogged()) {
        transition.redirect('login');
      }
    }
  },

  getInitialState() {
    return {
      userInfo: MachineStore.getUserInfo(),
      machineInfo: MachineStore.getMachineInfo(),
      activationInfo: MachineStore.getActivationInfo()
    };
  },

  getUserInfoToPassInProps() {
    var User = {
      Id: this.state.userInfo.Id,
      Role: this.state.userInfo.UserRole
    }
    return User;
  },

  /*
   * Logout with the exit button
   */
  handleLogout() {
    LoginActions.logout();
  },

  /*
   *
   */
  onChangeActivation() {
    this.setState({
      activationInfo: MachineStore.getActivationInfo()
    });
  },

  /*
   * To logout and redirect to login page
   */
  onChangeLogout() {
    if( !MachineStore.getIsLogged() ) {
      this.replaceWith('login');
    }
  },

  /*
   * Synchronize invent from store to machinepage
   */
  componentDidMount() {
    MachineStore.onChangeActivation = this.onChangeActivation;
    MachineStore.onChangeLogout = this.onChangeLogout;
  },

  render() {
    return (
      <div className="container-fluid" >
        <div>
          coucou {this.state.userInfo.FirstName} {this.state.userInfo.LastName}
        </div>
        <div className="container-fluid">
          <MachineList 
            user={this.getUserInfoToPassInProps()}
            info={this.state.machineInfo} 
            activation={this.state.activationInfo}
          />
        </div>
        <button 
          onClick={this.handleLogout}
          className="btn btn-danger" > Exit </button>
      </div>
    );
  }
});

module.exports = MachinePage;

import React from 'react'
import MachineList from './MachineList';
import MachineStore from '../stores/MachineStore';
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

  getUserId() {
    return this.state.userInfo.UserId;
  },

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
    console.log(this.state.activationInfo);
    return (
      <div className="container-fluid" >
        <div>
          coucou {this.state.userInfo.FirstName} {this.state.userInfo.LastName}
        </div>
        <div>
          <MachineList 
            uid={this.getUserId()}
            info={this.state.machineInfo} 
            activation={this.state.activationInfo}
          />
        </div>
      </div>
    );
  }
});

module.exports = MachinePage;

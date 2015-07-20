import $ from 'jquery';
import React from 'react'
import MachineList from './MachineList';
import MachineStore from '../stores/MachineStore';
import MachineActions from '../actions/MachineActions';
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
   * Create a table of the Id from an array
   * Used in shouldComponentUpdate to know get the id from previous state and next one
   */
  createCompareTable(state) {
    let table = [];
    for(let i in state) {
      table.push(state[i].Id);
    }
    return table;
  },

  shouldComponentUpdate(nextProps, nextState) {
    let shouldUpdate = false;
    let previousId = this.createCompareTable(this.state.activationInfo);
    let nextId = this.createCompareTable(nextState.activationInfo);
    shouldUpdate = $(previousId).not(nextId).length === 0 && $(nextId).not(previousId).length === 0;
    return !shouldUpdate;
  },

  /*
   * Destructor
   * Stop the polling
   */
  componentWillUnmount() {
    clearInterval(this.interval);
  },

  /*
   * Synchronize invent from store to machinepage
   */
  componentDidMount() {
    MachineStore.onChangeActivation = this.onChangeActivation;
    MachineStore.onChangeLogout = this.onChangeLogout;
    this.interval = setInterval(MachineActions.pollActivations, 1000);
  },

  render() {
    console.log('coucou');
    return (
      <div className="container-fluid" >
        <div>
          coucou {this.state.userInfo.FirstName} {this.state.userInfo.LastName}
        </div>
        <div >
          <MachineList 
            user={this.getUserInfoToPassInProps()}
            info={this.state.machineInfo} 
            activation={this.state.activationInfo}
          />
        </div>
        <button 
          onClick={this.handleLogout}
          className="btn btn-lg btn-block btn-danger" > Exit </button>
      </div>
    );
  }
});

module.exports = MachinePage;

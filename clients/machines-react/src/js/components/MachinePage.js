import $ from 'jquery';
import React from 'react';
import Flux from '../flux';
import getters from '../getters';
import MachineList from './MachineList';
import LoginStore from '../stores/LoginStore';
import MachineStore from '../stores/MachineStore';
import MachineActions from '../actions/MachineActions';
import LoginActions from '../actions/LoginActions';
import {Navigation} from 'react-router';

import toastr from 'toastr';
toastr.options.positionClass = 'toast-bottom-left';

/*
 * MachinePage:
 * Root component
 * Fetch the information from the store
 * Give it to its children to display the interface
 * TODO: reorganize and documente some function
 */
var MachinePage = React.createClass({

  /*
   * Enable some React router function as:
   *  ReplaceWith
   */
  mixins: [ Navigation, Flux.ReactMixin ],

  /*
   * If not logged then redirect to the login page
   */
  statics: {
    willTransitionTo(transition) {
      const isLogged = Flux.evaluateToJS(getters.getIsLogged);
      if(!isLogged) {
        transition.redirect('login');
      }
    }
  },

  /*
   * Start fetching the data
   * before the component is mounted
   */
  componentWillMount() {
    const uid = Flux.evaluateToJS(getters.getUid);
    MachineActions.fetchData(uid);
  },

  /*
   * Initial State
   * fetch data from MachineStore
   */
  getDataBindings() {
    return {
      userInfo: getters.getUserInfo,
      machineInfo: getters.getMachineInfo,
      activationInfo: getters.getActivationInfo
    };
  },

  /*
   * Callback called when nfc reader error occure
   */
  errorNFCCallback(error) {
    window.libnfc.cardRead.disconnect(this.nfcLogin);
    window.libnfc.cardReaderError.disconnect(this.errorNFCCallback);
    toastr.error(error);
    setTimeout(this.connectJsToQt, 2000);
  },

  /*
   * Return an object with information
   * Which are useful for MachineChooser
   */
  getUserInfoToPassInProps() {
    var User = {
      Id: this.state.userInfo.Id,
      Role: this.state.userInfo.UserRole
    };
    return User;
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

  /*
   * Look if the activations have a name
   * if they all have one, return true
   */
  hasNameInto(activation) {
    for(let i in activation ) {
      if(!activation.FirstName) {
        return false;
      }
    }
    return true;
  },

  /*
   * Clear state while logout
   */
  clearState() {
    MachineActions.clearState();
  },

  /*
   * Logout with the exit button
   */
  handleLogout() {
    LoginActions.logout();
  },

  /*
   * To logout and redirect to login page
   */
  onChangeLogout() {
    const isLogged = Flux.evaluateToJS(getters.getIsLogged);
    if (!isLogged) {
      this.replaceWith('login');
    }
  },

  /*
   * Destructor
   * Stop the polling
   */
  componentWillUnmount() {
    if(window.libnfc){
      window.libnfc.cardRead.disconnect(this.handleLogout);
      window.libnfc.cardReaderError.disconnect(this.errorNFCCallback);
    }
    this.clearState();
    clearInterval(this.interval);
  },

  /*
   * Call when the component is mounted in DOM
   * Synchronize invent from stores
   * Activate a polling (1,5s)
   */
  componentDidMount() {
    if(window.libnfc) {
      setTimeout(this.connectJsToQt, 1500);
    }
    MachineStore.onChangeActivation = this.onChangeActivation;
    LoginStore.onChangeLogout = this.onChangeLogout;
    MachineStore.onChangeLogin = this.onChangeLogin;
    this.interval = setInterval(MachineActions.pollActivations, 1500);

    Flux.observe(getters.getIsLogged, isLogged => {
      this.onChangeLogout();
    }.bind(this));
  },

  connectJsToQt() {
    toastr.info('You can now log out with your nfc card');
    window.libnfc.cardRead.connect(this.handleLogout);
    window.libnfc.cardReaderError.connect(this.errorNFCCallback);
    window.libnfc.asyncScan();
    setTimeout(this.handleLogout, 30000);
  },

  /*
   * Render the user name
   * MachinList
   * exit button
   */
  render() {
    var machineInfo = Flux.evaluateToJS(getters.getMachineInfo);
    return (
      <div>
        <div className="logged-user-name">
          <div className="text-center ng-binding">
            <i className="fa fa-user-secret"></i>&nbsp;
            {this.state.userInfo.FirstName} {this.state.userInfo.LastName}
          </div>
        </div>
        <MachineList
          user={this.getUserInfoToPassInProps()}
          info={machineInfo}
          activation={this.state.activationInfo}
        />
        <div className="container-fluid">
          <button
            onClick={this.handleLogout}
            className="btn btn-lg btn-block btn-danger btn-logout-bottom" > Exit </button>
        </div>
      </div>
    );
  }
});

module.exports = MachinePage;

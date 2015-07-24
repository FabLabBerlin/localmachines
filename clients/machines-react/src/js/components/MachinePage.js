import $ from 'jquery';
import React from 'react'
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
  mixins: [ Navigation ],

  /*
   * If not logged then redirect to the login page
   */
  statics: {
    willTransitionTo(transition) {
      if(!LoginStore.getIsLogged()) {
        transition.redirect('login');
      }
    }
  },

  /*
   * Start fetching the data 
   * before the component is mounted
   */
  componentWillMount() {
    let uid = LoginStore.getUid();
    MachineActions.fetchData(uid);
  },

  /*
   * Initial State
   * fetch data from MachineStore
   */
  getInitialState() {
    return {
      attemptToLog: true,
      userInfo: MachineStore.getUserInfo(),
      machineInfo: MachineStore.getMachineInfo(),
      activationInfo: MachineStore.getActivationInfo()
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
   * Fetch data to change the state after the initial render
   */
  fetchDataFromStore() {
    return {
      attemptToLog: false,
      userInfo: MachineStore.getUserInfo(),
      machineInfo: MachineStore.getMachineInfo(),
      activationInfo: MachineStore.getActivationInfo()
    };
  },

  /*
   * Return an object with information
   * Which are useful for MachineChooser
   */
  getUserInfoToPassInProps() {
    var User = {
      Id: this.state.userInfo.Id,
      Role: this.state.userInfo.UserRole
    }
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
    this.setState({
      attemptToLog: true,
      UserInfo: {},
      machineInfo: [],
      activationInfo: []
    });
    MachineActions.clearState();
  },

  /*
   * Logout with the exit button
   */
  handleLogout() {
    LoginActions.logout();
  },

  /*
   * When a new activation is fetch in the store
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
    if( !LoginStore.getIsLogged() ) {
      this.replaceWith('login');
    }
  },

  /*
   * When the fectching is done
   * Set state from the new data
   */
  onChangeLogin() {
    this.setState(this.fetchDataFromStore());
  },

  /*
   * Do not update (render) the component when false
   * If the data aren't not all loaded, return true
   * If some activation doesn't have a name loaded yet return true
   * Compare the activation id from the previous state with the new one
   * If they are the same, do not update
   *
   * WARNING: This function is really complicated BECAUSE the api call is badly
   * and the name is loaded in the store state asynchronously
   */
  shouldComponentUpdate(nextProps, nextState) {
    if( this.state.attemptToLog ) {
      return true;
    } else if(this.hasNameInto(this.state.activationInfo) === false){
      return true;
    } else {
      let shouldUpdate = false;
      let previousId = this.createCompareTable(this.state.activationInfo);
      let nextId = this.createCompareTable(nextState.activationInfo);
      shouldUpdate = $(previousId).not(nextId).length === 0 && $(nextId).not(previousId).length === 0;
      return !shouldUpdate;
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

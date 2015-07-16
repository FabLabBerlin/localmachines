
/*
 * Import toastr and set position
 */
import toastr from 'toastr';
toastr.options.positionClass = 'toast-bottom-left';

var MachineStore = {

  state: {
    isLogged: false,
    firstTry: true,
    userInfo: {},
    activationInfo:[],
    machineInfo: []
  },

  /*
   * POST call to the API
   * Make POST call cutomisable
   */
  postAPICall(url, dataToSend, successFunction,toastrMessage, errorFunction = function() {}) {
    $.ajax({
      url: url,
      dataType: 'json',
      type:'POST',
      data: dataToSend,
      success: function(data) {
        successFunction(data);
      }.bind(this),
      error: function(xhr, status, err) {
        if(toastrMessage != '') {
          toastr.error(toastrMessage);
        }
        errorFunction();
      }.bind(this),
    });
  },

  /*
   * GET call to the API
   * Make GET call cutomisable
   */
  getAPICall(url, successFunction, toastrMessage, errorFunction = function() {}) {
    $.ajax({
      url: url,
      dataType: 'json',
      type: 'GET',
      cache: false,
      success: function(data) {
        successFunction(data);
      }.bind(this),
      error: function(xhr, status, err) {
        if(toastrMessage != '') {
          toastr.error(toastrMessage);
        }
        errorFunction();
        console.log(url, status, err);
      }.bind(this),
    });
  },

  /*
   * API CALL TO DO
   * login
   *  - getname
   *  - get users/uid/machines
   *  - get activations/active
   *  - post activations
   *  - POST machines/mid/turn_onOrOff
   *  - put activation/aid
   *  - get name from user using machine
   * logout
   * 
   *
   */



  /*
   * Use POST call to login to the server
   * callback are defined below
   */
  postLogin(loginInfo){
    this.postAPICall('/api/users/login', loginInfo, this.LoginSuccess, '', this.LoginError);
  },

  /*
   * Use GET call to logout from the server
   * callback are defined below
   */
  getLogout() {
    this.getAPICall('/api/users/logout', this.logoutSuccess, '');
  },

  /*
   * Use GET call to get the name of the user
   * callback are defined below
   */
  getNameLogin(uid) {
    this.getAPICall('/api/users/' + uid + '/name', this.nameLoginSuccess, '');
  },

  /*
   * Use GET call to get the machines the user can use
   * callback are defined below
   */
  getUserMachines(uid) {
    this.getAPICall('/api/users/' + uid + '/machines', this.userMachineSuccess, '');
  },

  /*
   * Use GET call to get all the activations
   * callback are defined below
   */
  getActivationActive() {
    this.getAPICall('/api/activations/active', this.getActivationSuccess, '');
  },

  /*
   * To activate a machine
   * Create an activation
   */
  postActivation(mid){
    var dataToSend = {
      mid: mid
    };
    this.postAPICall('/api/activations', dataToSend, this.postActivationSuccess, 'Can not activate the machine');
  },


  /*
   *  Activated when getLogout succed
   *  MachineStore instead of this otherwise it doesn't work
   */
  logoutSuccess(data) {
    toastr.success('logout');
    MachineStore.state.isLogged = false;
    MachineStore.onChangeLogout();
  },

  /*
   * Activated when postLogin succeed
   * MachineStore instead of this otherwise it doesn't work
   */
  LoginSuccess(data) {
    var uid = data.UserId;
    MachineStore.state.userInfo.UserId = uid;
    MachineStore.state.firstTry = true;
    MachineStore.state.isLogged = true;
    MachineStore.getNameLogin(uid);
  },

  /*
   * Activated when getNameLogin succeed
   * MachineStore instead of this otherwe it doesn't work
   */
  nameLoginSuccess(data) {
    var uid = data.UserId;
    MachineStore.state.userInfo.FirstName = data.FirstName;
    MachineStore.state.userInfo.LastName = data.LastName;
    MachineStore.getUserMachines(uid)
  },

  /*
   * Activated when getUserMachines succeed
   */
  userMachineSuccess(data) {
    MachineStore.state.machineInfo = data;
    MachineStore.getActivationActive();
  },

  /*
   * Activated when getActivationActive succeed
   */
  getActivationSuccess(data) {
    var shortActivation = MachineStore.formatActivation(data);
    MachineStore.state.activationInfo = shortActivation;
    MachineStore.onChangeLogin();
  },

  formatActivation(rawActivation) {
    var shortActivation = [];
    var wantedInformation = ['Id', 'UserId', 'MachineId', 'TimeStart'];
    for( var i in rawActivation ) {
      var tmpItem = {};
      for( var indexWI in wantedInformation ) {
        tmpItem[wantedInformation[indexWI]] = rawActivation[i][wantedInformation[indexWI]];
      }
      shortActivation.push(tmpItem);
    }
    return shortActivation;
  },

  /*
   * Activated when login failed
   * MachineStore instead of this otherwise it doesn't work
   */
  LoginError(xhr, status, err) {
    if(MachineStore.state.firstTry === true) {
      MachineStore.state.firstTry = false;
    } else {
      toastr.error('Wrong password');
    }
  },

  endActivation(aid) {
    /*
    console.log('end activation');
    console.log(aid);
    */
  },

  postActivationSuccess(data) {
    toastr.success('activation done');
    MachineStore.onChangeActivation();
  },

  getIsLogged() {
    return this.state.isLogged;
  },

  getUserInfo() {
    return this.state.userInfo;
  },

  getActivationInfo() {
    return this.state.activationInfo;
  },

  getMachineInfo() {
    return this.state.machineInfo;
  },

  onChangeLogin() {},

  onChangeLogout() {},

  onChangeActivation() {}

}

module.exports = MachineStore;

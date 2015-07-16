
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
        console.error(url, status, err);
      }.bind(this),
    });
  },

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
   * Turn on a machine
   * Create an activation
   */
  postActivation(mid){
    var dataToSend = {
      mid: mid
    };
    this.postAPICall('/api/activations', dataToSend, this.postActivationSuccess, 'Can not activate the machine');
  },

  /*
   * Turn off a machine
   */
  putActivation(aid) {
    $.ajax({
      url: '/api/activations/' + aid,
      method: 'PUT',
      data: {
        ac: new Date().getTime()
      },
      success: function(data) {
        this.postActivationSuccess(data, 'Machine desactivated');
      }.bind(this),
      error: function(xhr, status, err){
        console.error('/api/activation/uid', status, err.toString());
      }.bind(this),
    });
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

  /*
   * Activated when postActivation succeed
   */
  postActivationSuccess(data, toastrMessage = 'Machine activated') {
    toastr.success(toastrMessage);
    var successFunction = function(data) {
      var shortActivation = MachineStore.formatActivation(data);
      MachineStore.state.activationInfo = shortActivation;
      MachineStore.onChangeActivation();
    }
    MachineStore.getAPICall('/api/activations/active', successFunction, '');
  },

  /*
   * Format rawActivation to have only useful information
   * @rawActivation: response send by the server
   */
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
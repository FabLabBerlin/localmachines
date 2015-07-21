
/*
 * Import toastr and set position
 */
import toastr from 'toastr';
toastr.options.positionClass = 'toast-bottom-left';

/*
 * Store the data
 * summary:
 * state (:34)
 * postAPI template function (:46)
 * getAPI template function (:69)
 * Call order (callback are define below or inside the apicall):
 *  - logout (:92)
 *  - login (:110)
 *  - userInfo (:128)
 *  - machineInfo (:151)
 *  - getActivationInfo (:167)
 *  - postActivationInfo (:190)
 *  - putActivation (:215)
 *  - postSwitchMachine (:238)
 * LoginError (:273)
 * utils:
 *  - formatActivation (:285)
 *  - nameInAllActivation (:301)
 *  - nameInOneActivation (:311)
 *  - getter (:324)
 *  - putLoginState (:343)
 *  - onChange (:353)
 */
var MachineStore = {

  /*
   * State of MachineStore
   * Information needed by the components
   */
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
  postAPICall(url, dataToSend, successFunction, toastrMessage = '', errorFunction = function() {}) {
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
        console.error(url, status, err);
      }.bind(this),
    });
  },

  /*
   * GET call to the API
   * Make GET call cutomisable
   */
  getAPICall(url, successFunction, toastrMessage = '', errorFunction = function() {}) {
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
   * Use GET call to logout from the server
   * callback are defined below
   */
  apiGetLogout() {
    this.getAPICall('/api/users/logout', this.logoutSuccess);
  },

  /*
   * Success Callback
   * Activated when getLogout succed
   * MachineStore instead of this otherwise it doesn't work
   */
  logoutSuccess(data) {
    toastr.success('logout');
    MachineStore.state.isLogged = false;
    MachineStore.onChangeLogout();
  },

  /*
   * Use POST call to login to the server
   * callback are defined below
   */
  apiPostLogin(loginInfo){
    this.postAPICall('/api/users/login', loginInfo, this.LoginSuccess, '', this.LoginError);
  },

  /*
   * Success Callback
   * Activated when postLogin succeed
   * MachineStore instead of this otherwise it doesn't work
   */
  LoginSuccess(data) {
    var uid = data.UserId;
    MachineStore.state.userInfo.UserId = uid;
    MachineStore.apiGetUserInfoLogin(uid);
  },

 /*
   * Use GET call to get the name of the user
   * callback are defined below
   */
  apiGetUserInfoLogin(uid) {
    this.getAPICall('/api/users/' + uid, this.userInfoSuccess);
  },

  /*
   * Success Callback
   * Activated when getNameLogin succeed
   * MachineStore instead of this otherwe it doesn't work
   */
  userInfoSuccess(data) {
    var uid = data.Id;
    var usefulInformation = ['Id', 'FirstName', 'LastName', 'UserRole'];
    var tmpInfo = {};
    for(var index in usefulInformation) {
      tmpInfo[usefulInformation[index]] = data[usefulInformation[index]];
    }
    MachineStore.state.userInfo = tmpInfo;
    MachineStore.apiGetUserMachines(uid)
  },

  /*
   * Use GET call to get the machines the user can use
   * callback are defined below
   */
  apiGetUserMachines(uid) {
    this.getAPICall('/api/users/' + uid + '/machines', this.userMachineSuccess);
  },

  /*
   * SuccessCallback
   * Activated when getUserMachines succeed
   */
  userMachineSuccess(data) {
    MachineStore.state.machineInfo = data;
    MachineStore.apiGetActivationActive();
  },

  /*
   * Use GET call to get all the activations
   * callback are defined below
   */
  apiGetActivationActive() {
    this.getAPICall('/api/activations/active', this.getActivationSuccess);
  },

  /*
   * Success Callback
   * Activated when getActivationActive succeed
   */
  getActivationSuccess(data) {
    var shortActivation = MachineStore.formatActivation(data);
    MachineStore.state.activationInfo = shortActivation;
    if(shortActivation.length != 0) {
      MachineStore.nameInAllActivations();
    } else if(MachineStore.state.isLogged === false){
      MachineStore.putLoginState();
    } else {
      MachineStore.onChangeActivation();
    }
  },

  /*
   * Turn on a machine
   * Create an activation
   */
  apiPostActivation(mid){
    var dataToSend = {
      mid: mid
    };
    this.postAPICall('/api/activations', dataToSend, this.postActivationSuccess, 'Can not activate the machine');
  },

  /*
   * Success Callback
   * Activated when postActivation succeed
   */
  postActivationSuccess(data, toastrMessage = 'Machine activated') {
    toastr.success(toastrMessage);
    var successFunction = function(data) {
      var shortActivation = MachineStore.formatActivation(data);
      MachineStore.state.activationInfo = shortActivation;
      MachineStore.onChangeActivation();
    }
    MachineStore.getAPICall('/api/activations/active', successFunction);
  },

  /*
   * End an activation
   * activation become unactive
   * @aid: activation id you want to shut down
   */
  apiPutActivation(aid) {
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
        console.error('/api/activation/aid', status, err.toString());
      }.bind(this),
    });
  },

  /*
   * Force a machine to be turned on or off
   * If the machine is active (activation) end the activation
   * @mid: machine you want to turn on or off
   * @onOrOff: action you want to do
   * @aid: activation id in case of turning off a machine
   * TODO: add animation
   */
  apiPostSwitchMachine(mid, onOrOff, aid = '') {
    //start animation
    if( onOrOff === 'off') {
      if(aid === '') {
        var successFunction = function(data) {
          //end animation
          toastr.success('machine off');
        };
      } else {
        var successFunction = function(data) {
          //end animation
          toastr.success('machine off and activation closed');
          MachineStore.apiPutActivation(aid);
        };
      }
    } else {
      var successFunction = function(data) {
        //end animation
        toastr.success('machine On');
      };
    }
    var errorFunction = function() {
      //end animation
    };
    this.postAPICall('/api/machines/' + mid + '/turn_' + onOrOff, 
                     { ac: new Date().getTime() },
                     successFunction,
                     'Fail to turn ' + onOrOff,
                     errorFunction
                    );
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

  /*
   * Format rawActivation to have only useful information
   * @rawActivation: response send by the server
   */
  formatActivation(rawActivation) {
    var shortActivation = [];
    var wantedInformation = ['Id', 'UserId', 'MachineId', 'TimeTotal'];
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
   * For each activation
   * Call nameInOneActivation
   */
  nameInAllActivations() {
    for( var index in this.state.activationInfo ) {
      var uid = this.state.activationInfo[index].UserId;
      this.nameInOneActivation(uid, index);
    }
  },

  /*
   * Put the name of the activation (identified by the index)
   * and put inside the json the name of the one who activates it
   */
  nameInOneActivation(uid, index) {
    var successFunction = function(data) {
      MachineStore.state.activationInfo[index]['FirstName'] = data.FirstName;
      MachineStore.state.activationInfo[index]['LastName'] = data.LastName;
      if(MachineStore.state.isLogged === false){
        MachineStore.putLoginState();
      } else {
        MachineStore.onChangeActivation();
      }
    };
    this.getAPICall('/api/users/' + uid + '/name', successFunction);
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

  /*
   * Change the state before login in
   */
  putLoginState() {
    this.state.isLogged = true;
    this.state.firstTry = true;
    this.onChangeLogin();
  },

  /*
   * Event triggered when login
   * See Login page
   */
  onChangeLogin() {},

  /*
   * Event triggered when logout
   * See MachinePage
   */
  onChangeLogout() {},

  /*
   * Event triggered when activation is send or get
   * See MachinePage
   */
  onChangeActivation() {}

}

module.exports = MachineStore;

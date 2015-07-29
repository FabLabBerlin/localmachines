import $ from 'jquery';
import _ from 'lodash';

/*
 * import toastr and set position
 */
import toastr from 'toastr';
toastr.options.positionClass = 'toast-bottom-left';

/*
 * TODO: refactoring some comments and function
 * Store the data
 * summary:
 * state (:34)
 * postAPI template function (:46)
 * getAPI template function (:69)
 * Call order (callback are define below or inside the apicall):
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
    activationInfo: [],
    machineInfo: []
  },

  /*
   * POST call to the API
   * Make POST call cutomisable
   */
  _postAPICall(url, dataToSend, successFunction, toastrMessage = '', errorFunction = function() {}) {
    $.ajax({
      url: url,
      dataType: 'json',
      type: 'POST',
      data: dataToSend,
      success: function(data) {
        successFunction(data);
      },
      error: function(xhr, status, err) {
        if (toastrMessage) {
          toastr.error(toastrMessage);
        }
        errorFunction();
        console.error(url, status, err);
      }
    });
  },

  /*
   * GET call to the API
   * Make GET call cutomisable
   */
  _getAPICall(url, successFunction, toastrMessage = '', errorFunction = function() {}) {
    $.ajax({
      url: url,
      dataType: 'json',
      type: 'GET',
      cache: false,
      success: function(data) {
        successFunction(data);
      },
      error: function(xhr, status, err) {
        if (toastrMessage) {
          toastr.error(toastrMessage);
        }
        errorFunction();
        console.error(url, status, err);
      }
    });
  },

 /*
   * Use GET call to get the name of the user
   * callback are defined below
   */
  apiGetUserInfoLogin(uid) {
    this._getAPICall('/api/users/' + uid, this._userInfoSuccess);
  },

  /*
   * Success Callback
   * Activated when getNameLogin succeed
   * MachineStore instead of this otherwe it doesn't work
   */
  _userInfoSuccess(data) {
    var uid = data.Id;
    var usefulInformation = ['Id', 'FirstName', 'LastName', 'UserRole'];
    var tmpInfo = {};
    for(var index in usefulInformation) {
      tmpInfo[usefulInformation[index]] = data[usefulInformation[index]];
    }
    MachineStore.state.userInfo = tmpInfo;
    MachineStore.apiGetUserMachines(uid);
  },

  /*
   * Use GET call to get the machines the user can use
   * callback are defined below
   */
  apiGetUserMachines(uid) {
    this._getAPICall('/api/users/' + uid + '/machines', this._userMachineSuccess);
  },

  /*
   * SuccessCallback
   * Activated when getUserMachines succeed
   */
  _userMachineSuccess(data) {
    MachineStore.state.machineInfo = data;
    MachineStore.apiGetActivationActive();
  },

  /*
   * Use GET call to get all the activations
   * callback are defined below
   */
  apiGetActivationActive() {
    this._getAPICall('/api/activations/active', this._getActivationSuccess);
  },

  /*
   * Success Callback
   * Activated when getActivationActive succeed
   */
  _getActivationSuccess(data) {
    var shortActivation = MachineStore._formatActivation(data);
    MachineStore.state.activationInfo = shortActivation;
    if (shortActivation.length !== 0) {
      MachineStore._nameInAllActivations();
    } else if (!MachineStore.state.isLogged){
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
    this._postAPICall('/api/activations', dataToSend, this._postActivationSuccess, 'Can not activate the machine');
  },

  /*
   * Success Callback
   * Activated when postActivation succeed
   */
  _postActivationSuccess(data, toastrMessage = 'Machine activated') {
    toastr.success(toastrMessage);
    var successFunction = function(getData) {
      var shortActivation = MachineStore._formatActivation(getData);
      MachineStore.state.activationInfo = shortActivation;
      MachineStore.onChangeActivation();
    };
    MachineStore._getAPICall('/api/activations/active', successFunction);
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
        this._postActivationSuccess(data, 'Machine desactivated');
      }.bind(this),
      error: function(xhr, status, err) {
        console.error('/api/activation/aid', status, err.toString());
      }.bind(this)
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
    var successFunction;
    //start animation
    if( onOrOff === 'off') {
      if(aid === '') {
        successFunction = function(data) {
          //end animation
          toastr.success('machine off');
        };
      } else {
        successFunction = function(data) {
          //end animation
          toastr.success('machine off and activation closed');
          MachineStore.apiPutActivation(aid);
        };
      }
    } else {
      successFunction = function(data) {
        //end animation
        toastr.success('machine On');
      };
    }
    var errorFunction = function() {
      //end animation
    };
    this._postAPICall('/api/machines/' + mid + '/turn_' + onOrOff,
                     { ac: new Date().getTime() },
                     successFunction,
                     'Fail to turn ' + onOrOff,
                     errorFunction
                    );
  },

  /*
   * Format rawActivation to have only useful information
   * @rawActivation: response send by the server
   */
  _formatActivation(rawActivation) {
    return _.map(rawActivation, function(rawActivationItem) {
      var tmpItem = {};
      ['Id', 'UserId', 'MachineId', 'TimeTotal'].forEach(function(key){
        tmpItem[key] = rawActivationItem[key];
      });
      return tmpItem;
    });
  },

  /*
   * For each activation
   * Call nameInOneActivation
   */
  _nameInAllActivations() {
    this.state.activationInfo.forEach(function(activation, i) {
      this._nameInOneActivation(activation, i);
    }.bind(this));
  },

  /*
   * Put the name of the activation (identified by the index)
   * and put inside the json the name of the one who activates it
   */
  _nameInOneActivation(activation, index) {
    this._getAPICall('/api/users/' + activation.UserId + '/name', function(userData) {
      _.merge(activation, userData);
      if (!MachineStore.state.isLogged) {
        MachineStore.putLoginState();
      } else {
        MachineStore.onChangeActivation();
      }
    });
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
   * Clean State before login out
   */
  clearState() {
    this.state = {
      isLogged: false,
      firstTry: true,
      userInfo: {},
      activationInfo: [],
      machineInfo: []
    };
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

};

module.exports = MachineStore;

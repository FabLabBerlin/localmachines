var $ = require('jquery');
var actionTypes = require('../actionTypes');
var getters = require('../getters');
var MachineStore = require('../stores/MachineStore');
var reactor = require('../reactor');
var toastr = require('../toastr');


var MachineActions = {

  fetchData(uid) {
    _getAPICall('/api/users/' + uid, _userInfoSuccess);
  },

  /*
   * To end an activation
   * @aid: id of the activation you want to shut down
   */
  endActivation(aid) {
    apiPutActivation(aid);
  },

  /*
   * To start an activation
   * @mid: id of the machine you want to activate
   */
  startActivation(mid) {
    var dataToSend = {
      mid: mid
    };
    _postAPICall('/api/activations', dataToSend, _postActivationSuccess, 'Can not activate the machine');
  },

  /*
   * When an admin want to force on a machine
   */
  adminTurnOffMachine(mid, aid) {
    apiPostSwitchMachine(mid, 'off', aid);
  },

  /*
   * When an admin want to force off a machine
   */
  adminTurnOnMachine(mid) {
    apiPostSwitchMachine(mid, 'on');
  },

  /*
   * Clear store state while logout
   */
  clearState() {
    reactor.dispatch(actionTypes.MACHINE_STORE_CLEAR_STATE);
  },

  /*
   * To continue to refresh the view each seconds
   */
  pollActivations() {
    _getAPICall('/api/activations/active', function(data) {
      var activationInfo = _formatActivation(data);
      reactor.dispatch(actionTypes.SET_ACTIVATION_INFO, { activationInfo });
    });
  }

};

/*
 * Use GET call to get the name of the user
 * callback are defined below
 */
function apiGetUserInfoLogin(uid) {
  _getAPICall('/api/users/' + uid, _userInfoSuccess);
}

/*
 * Use GET call to get the machines the user can use
 * callback are defined below
 */
function apiGetUserMachines(uid) {
  _getAPICall('/api/users/' + uid + '/machines', function(machineInfo) {
    reactor.dispatch(actionTypes.SET_MACHINE_INFO, { machineInfo });
    apiGetActivationActive();
  });
}

/*
 * End an activation
 * activation become unactive
 * @aid: activation id you want to shut down
 */
function apiPutActivation(aid) {
  $.ajax({
    url: '/api/activations/' + aid,
    method: 'PUT',
    data: {
      ac: new Date().getTime()
    },
    success: function(data) {
      _postActivationSuccess(data, 'Machine deactivated');
    }.bind(this),
    error: function(xhr, status, err) {
      toastr.error('Failed to deactivate');
      console.error('/api/activation/aid', status, err.toString());
    }.bind(this)
  });
}

/*
 * Force a machine to be turned on or off
 * If the machine is active (activation) end the activation
 * @mid: machine you want to turn on or off
 * @onOrOff: action you want to do
 * @aid: activation id in case of turning off a machine
 * TODO: add animation
 */
function apiPostSwitchMachine(mid, onOrOff, aid = '') {
  var successFunction;
  //start animation
  if( onOrOff === 'off') {
    if(aid === '') {
      successFunction = function(data) {
        //end animation
        toastr.success('Machine off');
      };
    } else {
      successFunction = function(data) {
        //end animation
        toastr.success('Machine off and activation closed');
        apiPutActivation(aid);
      };
    }
  } else {
    successFunction = function(data) {
      //end animation
      toastr.success('Machine on');
    };
  }
  var errorFunction = function() {
    //end animation
  };
  _postAPICall('/api/machines/' + mid + '/turn_' + onOrOff,
                   { ac: new Date().getTime() },
                   successFunction,
                   'Failed to turn ' + onOrOff,
                   errorFunction
                  );
}

function apiGetActivationActive() {
  _getAPICall('/api/activations/active', function(data) {
    var activationInfo = _formatActivation(data);
    _.each(activationInfo, apiLoadMachineUser);
    reactor.dispatch(actionTypes.SET_ACTIVATION_INFO, { activationInfo });
  });
}

function apiLoadMachineUser(activation) {
  _getAPICall('/api/users/' + activation.UserId + '/name', function(userData) {
    reactor.dispatch(actionTypes.REGISTER_MACHINE_USER, { userData });
  });
}

/*
 * GET call to the API
 * Make GET call cutomisable
 */
function _getAPICall(url, successFunction, toastrMessage = '', errorFunction = function() {}) {
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
}

/*
 * POST call to the API
 * Make POST call cutomisable
 */
function _postAPICall(url, dataToSend, successFunction, toastrMessage = '', errorFunction = function() {}) {
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
}

/*
 * Format rawActivation to have only useful information
 * @rawActivation: response send by the server
 */
function _formatActivation(rawActivation) {
  return _.map(rawActivation, function(rawActivationItem) {
    var tmpItem = {};
    ['Id', 'UserId', 'MachineId', 'TimeTotal'].forEach(function(key){
      tmpItem[key] = rawActivationItem[key];
    });
    return tmpItem;
  });
}

/*
 * Success Callback
 * Activated when postActivation succeed
 */
function _postActivationSuccess(data, toastrMessage = 'Machine activated') {
  toastr.success(toastrMessage);
  var successFunction = function(getData) {
    var activationInfo = _formatActivation(getData);
    reactor.dispatch(actionTypes.SET_ACTIVATION_INFO, { activationInfo });
    _.each(activationInfo, apiLoadMachineUser);
  };
  _getAPICall('/api/activations/active', successFunction);
}

/*
 * Success Callback
 * Activated when getNameLogin succeed
 * MachineStore instead of this otherwe it doesn't work
 */
function _userInfoSuccess(data) {
  var uid = data.Id;
  var usefulInformation = ['Id', 'FirstName', 'LastName', 'UserRole'];
  var userInfo = {};
  for(var index in usefulInformation) {
    userInfo[usefulInformation[index]] = data[usefulInformation[index]];
  }
  reactor.dispatch(actionTypes.SET_USER_INFO, { userInfo });
  apiGetUserMachines(uid);
}

export default MachineActions;

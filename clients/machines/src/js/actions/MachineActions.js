var $ = require('jquery');
var actionTypes = require('../actionTypes');
var ApiActions = require('./ApiActions');
var getters = require('../getters');
var GlobalActions = require('./GlobalActions');
var LoginActions = require('../actions/LoginActions');
var MachineStore = require('../stores/MachineStore');
var reactor = require('../reactor');
var toastr = require('../toastr');


var MachineActions = {

  /*
   * To end an activation
   * @aid: id of the activation you want to shut down
   */
  endActivation(aid, cb) {
    apiPutActivation(aid, cb);
    LoginActions.keepAlive();
  },

  /*
   * To start an activation
   * @mid: id of the machine you want to activate
   */
  startActivation(mid) {
    var dataToSend = {
      mid: mid
    };
    ApiActions.postCall('/api/activations', dataToSend, _postActivationSuccess, 'Can not activate the machine');
    LoginActions.keepAlive();
  },

  /*
   * When an admin want to force on a machine
   */
  adminTurnOffMachine(mid, aid) {
    apiPostSwitchMachine(mid, 'off', aid);
    LoginActions.keepAlive();
  },

  /*
   * When an admin want to force off a machine
   */
  adminTurnOnMachine(mid) {
    apiPostSwitchMachine(mid, 'on');
    LoginActions.keepAlive();
  },

  /*
   * Use GET call to get the machines the user can use
   * callback are defined below
   */
  apiGetUserMachines(uid) {
    ApiActions.getCall('/api/users/' + uid + '/machines', function(machineInfo) {
      reactor.dispatch(actionTypes.SET_MACHINE_INFO, { machineInfo });
      apiGetActivationActive();
    });
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
    ApiActions.getCall('/api/activations/active', function(data) {
      var activationInfo = _formatActivation(data);
      reactor.dispatch(actionTypes.SET_ACTIVATION_INFO, { activationInfo });
    });
  },

  pollMachines() {
    const uid = reactor.evaluateToJS(getters.getUid);
    ApiActions.getCall('/api/users/' + uid + '/machines', function(machineInfo) {
      reactor.dispatch(actionTypes.SET_MACHINE_INFO, { machineInfo });
    });
  },

  setUnderMaintenance({ mid, onOrOff }) {
    ApiActions.postCall('/api/machines/' + mid + '/under_maintenance/' + onOrOff,
                    {},
                    function(data) {
                      reactor.dispatch(actionTypes.SET_UNDER_MAINTENANCE, { mid, onOrOff });
                      if (onOrOff === 'on') {
                        toastr.info('Machine under maintenance');
                      } else {
                        toastr.info('Machine is working again');
                      }
                    },
                    function() {
                      toastr.error('Could not change maintenance mode');
                    }
    );
  }

};

/*
 * End an activation
 * activation become unactive
 * @aid: activation id you want to shut down
 */
function apiPutActivation(aid, cb) {
  GlobalActions.showGlobalLoader();
  $.ajax({
    url: '/api/activations/' + aid,
    method: 'PUT',
    data: {
      ac: new Date().getTime()
    },
    success: function(data) {
      GlobalActions.hideGlobalLoader();
      _postActivationSuccess(data, 'Machine deactivated');
      if (cb) {
        cb();
      }
    }.bind(this),
    error: function(xhr, status, err) {
      GlobalActions.hideGlobalLoader();
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
  if( onOrOff === 'off') {
    if(aid === '') {
      successFunction = function(data) {
        toastr.success('Machine off');
      };
    } else {
      successFunction = function(data) {
        toastr.success('Machine off and activation closed');
        apiPutActivation(aid);
      };
    }
  } else {
    successFunction = function(data) {
      toastr.success('Machine on');
    };
  }
  var errorFunction = function() {
  };
  ApiActions.postCall('/api/machines/' + mid + '/turn_' + onOrOff,
                   { ac: new Date().getTime() },
                   successFunction,
                   'Failed to turn ' + onOrOff,
                   errorFunction
                  );
}

function apiGetActivationActive() {
  ApiActions.getCall('/api/activations/active', function(data) {
    var activationInfo = _formatActivation(data);
    _.each(activationInfo, apiLoadMachineUser);
    reactor.dispatch(actionTypes.SET_ACTIVATION_INFO, { activationInfo });
  });
}

function apiLoadMachineUser(activation) {
  ApiActions.getCall('/api/users/' + activation.UserId + '/name', function(userData) {
    reactor.dispatch(actionTypes.REGISTER_MACHINE_USER, { userData });
  });
}

/*
 * Format rawActivation to have only useful information
 * @rawActivation: response send by the server
 */
function _formatActivation(rawActivation) {
  return _.map(rawActivation, function(rawActivationItem) {
    var tmpItem = {};
    ['Id', 'UserId', 'MachineId', 'Quantity'].forEach(function(key){
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
  ApiActions.getCall('/api/activations/active', successFunction);
}

export default MachineActions;

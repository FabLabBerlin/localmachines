var $ = require('jquery');
var actionTypes = require('./actionTypes');
var ApiActions = require('../Api/actions');
var getters = require('./getters');
var GlobalActions = require('../Global/actions');
var LoginActions = require('../Login/actions');
var MachineStore = require('./stores/MachineStore');
var reactor = require('../../reactor');
var toastr = require('../../toastr');


var MachineActions = {

  startActivation(mid) {
    var dataToSend = {
      mid: mid
    };
    GlobalActions.showGlobalLoader();
    $.ajax({
      url: '/api/activations',
      dataType: 'json',
      type: 'POST',
      data: {
        mid: mid
      },
      success(data) {
        GlobalActions.hideGlobalLoader();
        toastr.info('Machine activated');
      },
      error(xhr, status, err) {
        GlobalActions.hideGlobalLoader();
        toastr.error('Can not activate the machine');
        console.error(status, err);
      }
    });
    LoginActions.keepAlive();
  },

  endActivation(aid, cb) {
    endActivation(aid, cb);
    LoginActions.keepAlive();
  },

  forceTurnOnMachine(mid) {
    GlobalActions.showGlobalLoader();
    $.ajax({
      url: '/api/machines/' + mid + '/turn_on',
      type: 'POST',
      success(data) {
        GlobalActions.hideGlobalLoader();
        toastr.success('Machine on');
      },
      error(xhr, status, err) {
        GlobalActions.hideGlobalLoader();
        toastr.error('Failed to turn on');
        console.error(status, err);
      }
    });
    LoginActions.keepAlive();
  },

  forceTurnOffMachine(mid, aid) {
    GlobalActions.showGlobalLoader();
    $.ajax({
      url: '/api/machines/' + mid + '/turn_off',
      type: 'POST',
      success(data) {
        GlobalActions.hideGlobalLoader();
        if (aid) {
          toastr.success('Machine off and activation closed');
          endActivation(aid);
        } else {
          toastr.success('Machine off');
        }
      },
      error(xhr, status, err) {
        GlobalActions.hideGlobalLoader();
        toastr.error('Failed to turn off');
        console.error(status, err);
      }
    });
    LoginActions.keepAlive();
  },

  /*
   * Use GET call to get the machines the user can use
   * callback are defined below
   */
  apiGetUserMachines(uid) {
    ApiActions.getCall('/api/users/' + uid + '/machines', (machines) => {
      reactor.dispatch(actionTypes.SET_MACHINES, { machines });
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
  pollDashboard() {
    const uid = reactor.evaluateToJS(getters.getUid);
    ApiActions.getCall('/api/users/' + uid + '/dashboard', function(data) {
      var userIds = [];

      reactor.dispatch(actionTypes.SET_ACTIVATIONS, {
        activations: data.Activations
      });
      if (data.Activations) {
        userIds = _.pluck(data.Activations, 'UserId');
      }
      reactor.dispatch(actionTypes.SET_MACHINES, {
        machines: data.Machines
      });
      if (data.Tutorings) {
        reactor.dispatch(actionTypes.SET_TUTORINGS, data.Tutorings.Data);
        userIds = _.union(userIds, _.pluck(data.Tutorings.Data, 'UserId'));
        userIds = _.filter(userIds, (uid) => {
          return uid;
        });
      }
      fetchUserNames(userIds);
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
function endActivation(aid, cb) {
  GlobalActions.showGlobalLoader();
  $.ajax({
    url: '/api/activations/' + aid,
    method: 'PUT',
    data: {
      ac: new Date().getTime()
    },
    success(data) {
      GlobalActions.hideGlobalLoader();
      toastr.info('Machine deactivated');
      if (cb) {
        cb();
      }
    },
    error(xhr, status, err) {
      GlobalActions.hideGlobalLoader();
      toastr.error('Failed to deactivate');
      console.error('/api/activation/aid', status, err.toString());
    }
  });
}

function fetchUserNames(userIds) {
  var fetchedUserIds = _.keys(reactor.evaluateToJS(getters.getMachineUsers));
  fetchedUserIds = _.map(fetchedUserIds, function(id) {
    return parseInt(id, 10);
  });
  userIds = _.difference(userIds, fetchedUserIds);

  if (userIds.length > 0) {
    $.ajax({
      url: '/api/users/names?uids=' + userIds.join(','),
      dataType: 'json',
      type: 'GET',
      success(response) {
        _.each(response.Users, function(userData) {
          reactor.dispatch(actionTypes.REGISTER_MACHINE_USER, { userData });
        });
      },
      error() {
        console.log('Error loading names');
      }
    });
  }
}

export default MachineActions;

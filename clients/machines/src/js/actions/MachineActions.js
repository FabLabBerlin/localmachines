var _ = require('lodash');
var $ = require('jquery');
var actionTypes = require('../actionTypes');
var ApiActions = require('./ApiActions');
var getters = require('../getters');
var GlobalActions = require('./GlobalActions');
var LocationGetters = require('../modules/Location/getters');
var LoginActions = require('../actions/LoginActions');
var MachineStore = require('../stores/MachineStore');
var reactor = require('../reactor');
var toastr = require('../toastr');


var MachineActions = {

  startActivation(mid) {
    var dataToSend = {
      mid: mid
    };
    var locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    GlobalActions.showGlobalLoader();
    $.ajax({
      url: '/api/activations/start?location=' + locationId,
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
        if (xhr.status === 403 && status === 'No remote activation') {
          toastr.error('Activations only possible through Lab Wifi for safety reasons');
        } else {
          toastr.error('Can not activate the machine');
          console.error(status, err);
        }
      }
    });
    LoginActions.keepAlive();
  },

  endActivation(aid, cb) {
    endActivation(aid, cb);
    LoginActions.keepAlive();
  },

  /*
   * Use GET call to get the machines the user can use
   * callback are defined below
   */
  apiGetUserMachines(locationId, uid) {
    var url = '/api/users/' + uid + '/machines?location=' + locationId;
    ApiActions.getCall(url, (machines) => {
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
  pollDashboard(router, locationId) {
    const uid = reactor.evaluateToJS(getters.getUid);

    $.ajax({
      url: '/api/users/' + uid + '/dashboard?location=' + locationId,
      dataType: 'json',
      type: 'GET',
      cache: false,
      success: function(data) {
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
          userIds = _.filter(userIds, (userId) => {
            return userId;
          });
        }
        fetchUserNames(userIds);
      },
      error: function(xhr, status) {
        console.log('xhr:', xhr);
        if (xhr.status === 401 && xhr.responseText === 'Not logged in') {
          toastr.error('Session not active anymore. Logging out.');
          LoginActions.logout(router);
        }
      }
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
    url: '/api/activations/' + aid + '/close',
    method: 'POST',
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

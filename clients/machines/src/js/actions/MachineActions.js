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


var lpRequest;
var socket;

function dashboardDispatch(data) {
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
}


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

  lpDashboard(router, locationId, chained) {
    if (lpRequest && !chained) {
      lpRequest.abort();
      lpRequest = undefined;
    }

    const uid = reactor.evaluateToJS(getters.getUid);
    var url;
    if (chained) {
      url = '/api/users/' + uid + '/dashboard/lp?location=' + locationId;
    } else {
      url = '/api/users/' + uid + '/dashboard?location=' + locationId;
    }

    lpRequest = $.ajax({
      url: url,
      dataType: 'json',
      type: 'GET',
      cache: false,
      timeout: 30000
    });

    lpRequest.then(function(data) {
      dashboardDispatch(data);
      MachineActions.lpDashboard(router, locationId, true);
    });
    lpRequest.fail(function(xhr, status) {
      if (status === 'abort') {
        console.log('abort lp');
        return;
      }

      console.log('xhr:', xhr);
      if (xhr.status === 401 && xhr.responseText === 'Not logged in') {
        toastr.error('Session not active anymore. Logging out.');
        LoginActions.logout(router);
      } else {
        toastr.error('Connection error.  Reconnecting...');
        console.log('reconnecting in 5 s...');
        window.setTimeout(function() {
          MachineActions.lpDashboard(router, locationId);
        }, 5000);
      }
    });
  },

  wsDashboard(router, locationId) {
    const t0 = new Date();
    if (lpRequest) {
      lpRequest.abort();
      lpRequest = undefined;
    }
    if (socket) {
      socket.onclose = function () {}; // disable onclose handler first
      socket.close();
      socket = undefined;
    }

    const uid = reactor.evaluateToJS(getters.getUid);
    const host = window.location.host;
    const protocol = host === 'easylab.io' ? 'wss' : 'ws';
    socket = new WebSocket(protocol + '://' + host + '/api/users/' + uid + '/dashboard/ws?location=' + locationId);
    socket.onmessage = function(e) {
      dashboardDispatch(JSON.parse(e.data));
    };
    socket.onclose = function(e) {
      var duration = new Date() - t0;
      if (duration < 30000) {
        toastr.warning('Falling back to longpoll...');
        MachineActions.lpDashboard(router, locationId);
      } else {
        socket = null;
        console.log('reconnecting in 5 s...');
        window.setTimeout(function() {
          MachineActions.wsDashboard(router, locationId);
        }, 5000);
      }
    };
    socket.onerror = function(e) {
      console.log('websocket error:', e);
    };
  },

  pollDashboard() {},

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
        reactor.dispatch(actionTypes.REGISTER_MACHINE_USERS, response.Users);
      },
      error() {
        console.log('Error loading names');
      }
    });
  }
}

export default MachineActions;

import _ from 'lodash';
var $ = require('jquery');
import actionTypes from '../actionTypes';
import getters from '../getters';
import GlobalActions from './GlobalActions';
import LocationGetters from '../modules/Location/getters';
import LoginActions from '../actions/LoginActions';
import Machines from '../modules/Machines';
import reactor from '../reactor';
import toastr from '../toastr';


var lpRequest;
var socket;

function dashboardDispatch(data) {
  var userIds = [];

  if (data.UserMessage && data.UserMessage.Error) {
    toastr.error(data.UserMessage.Error);
  }
  if (data.UserMessage && data.UserMessage.Info) {
    toastr.info(data.UserMessage.Info);
  }
  if (data.UserMessage && data.UserMessage.Warning) {
    toastr.warn(data.UserMessage.Warning);
  }
  if (data.Activations || data.Machines) {
    reactor.batch(() => {
      reactor.dispatch(Machines.actionTypes.SET_ACTIVATIONS, {
        activations: data.Activations
      });
      if (data.Activations) {
        userIds = _.map(data.Activations, 'UserId');
      }
      reactor.dispatch(Machines.actionTypes.SET_MACHINES, {
        machines: data.Machines
      });
      MachineActions.fetchUserNames(userIds);
    });
  }
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
        window.setTimeout(() => {
          GlobalActions.hideGlobalLoader();
        }, 1000);
        toastr.info('Sent On Request');
      },
      error(xhr, status, err) {
        if (xhr.responseText === 'Not logged in') {
          window.location.href = '/logout';
        } else {
          GlobalActions.hideGlobalLoader();
          if (xhr.status === 403 && xhr.responseText === 'No remote activation') {
            toastr.error('Must start in Lab Wifi for safety (VPN won\'t work either)');
          } else {
            toastr.error('Cannot request turn on');
            console.error(status, err);
          }
        }
      }
    });
  },

  endActivation(aid, cb) {
    var locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    GlobalActions.showGlobalLoader();
    $.ajax({
      url: '/api/activations/' + aid + '/close?location=' + locationId,
      method: 'POST',
      data: {
        ac: new Date().getTime()
      },
      success(data) {
        window.setTimeout(() => {
          GlobalActions.hideGlobalLoader();
        }, 1000);
        toastr.info('Sent Off Request');
        if (cb) {
          cb();
        }
      },
      error(xhr, status, err) {
        if (xhr.responseText === 'Not logged in') {
          window.location.href = '/logout';
        } else {
          GlobalActions.hideGlobalLoader();
          toastr.error('Failed to deactivate');
          console.error('/api/activation/aid', status, err.toString());
        }
      }
    });
  },

  apiGetUserMachines(locationId, uid) {
    $.ajax({
      url: '/api/users/' + uid + '/machines?location=' + locationId
    })
    .done(machines => {
      reactor.dispatch(Machines.actionTypes.SET_MACHINES, { machines });
    });
  },

  /*
   * Clear store state while logout
   */
  clearState() {
    reactor.dispatch(Machines.actionTypes.MACHINE_STORE_CLEAR_STATE);
  },

  fetchUserNames(userIds) {
    var fetchedUserIds = _.keys(reactor.evaluateToJS(Machines.getters.getMachineUsers));
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
          reactor.dispatch(Machines.actionTypes.REGISTER_MACHINE_USERS, response.Users);
        },
        error() {
          console.log('Error loading names');
        }
      });
    }
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
    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
    socket = new WebSocket(protocol + '://' + host + '/api/users/' + uid + '/dashboard/ws?location=' + locationId);
    socket.onmessage = function(e) {
      dashboardDispatch(JSON.parse(e.data));
    };
    socket.onclose = function(e) {
      console.log('WebSocket closed, e=', e);
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
    $.ajax({
      url: '/api/machines/' + mid + '/under_maintenance/' + onOrOff,
      type: 'POST'
    })
    .done(data => {
      reactor.dispatch(Machines.actionTypes.SET_UNDER_MAINTENANCE, { mid, onOrOff });
      if (onOrOff === 'on') {
        toastr.info('Machine under maintenance');
      } else {
        toastr.info('Machine is working again');
      }
    })
    .error(() => {
      toastr.error('Could not change maintenance mode');
    });
  },

  updateMachineField(mid, name, value) {
    reactor.dispatch(Machines.actionTypes.UPDATE_MACHINE_FIELD, {mid, name, value});
  },

  uploadMachineImage(mid, data) {
    reactor.dispatch(Machines.actionTypes.UPLOAD_MACHINE_IMAGE, {mid, data});
  }

};

/*
 * End an activation
 * activation become unactive
 * @aid: activation id you want to shut down
 */
function endActivation(aid, cb) {

}

export default MachineActions;

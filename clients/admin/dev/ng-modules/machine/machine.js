(function(){

'use strict';

var app = angular.module('fabsmith.admin.machine', 
 ['ngRoute', 'ngCookies', 'fabsmith.admin.randomtoken']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/machine/:machineId', {
    templateUrl: 'ng-modules/machine/machine.html',
    controller: 'MachineCtrl'
  });
}]); // app.config

app.controller('MachineCtrl', 
 ['$rootScope', '$scope', '$cookies', '$routeParams', '$http', '$location', '$filter', '$timeout', 'randomToken', 'api',
 function($rootScope, $scope, $cookies, $routeParams, $http, $location, $filter, $timeout, randomToken, api) {

  $scope.mainMenu = $rootScope.mainMenu;

  $scope.machine = {
    Id: $routeParams.machineId
  };

  $scope.machineImageFile = undefined;
  $scope.machineImageNewFile = undefined;
  $scope.machineImageNewFileName = undefined;
  $scope.machineImageNewFileSize = undefined;
  $scope.netswitchConfigStatus = undefined;
  $scope.unsavedChanges = false;
  $scope.loading = false;
  $scope.user = undefined;

  $scope.registerUnsavedChange = function() {
    $scope.unsavedChanges = true;
  };

  $scope.loadMachine = function() {
    var machineId = $scope.machine.Id;
    $http({
      method: 'GET',
      url: '/api/machines/' + machineId,
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      $scope.machine = data;
      $scope.machine.TypeId = String($scope.machine.TypeId);
      $scope.machine.Price = $filter('currency')($scope.machine.Price,'',2);
      if ($scope.machine.Image) {
        $scope.machineImageFile = "/files/" + $scope.machine.Image;
      }
    })
    .error(function() {
      toastr.error('Failed to get machine');
    });
  };

  $scope.loadMachineTypes = function() {
    $http({
      method: 'GET',
      url: '/api/machine_types'
    })
    .success(function(data) {
      $scope.machineTypes = data;
    })
    .error(function() {
      toastr.error('Failed to get machine types');
    });
  };

  $scope.loadMachineTypes();
  $scope.loadMachine();

  $scope.setArchived = function(archived) {
    var action = archived ? 'archived' : 'unarchived';

    $http({
      method: 'POST',
      url: '/api/machines/' + $scope.machine.Id + '/set_archived',
      params: {
        archived: archived
      }
    })
    .success(function(data) {
      toastr.info('Successfully ' + action + ' machine');
      $scope.loadMachine();
    })
    .error(function() {
      toastr.error('Failed to ' + action + ' machine');
    });
  }; // archive()

  $scope.updateMachine = function() {
    $scope.unsavedChanges = false;

    // Make a clone of the machine model
    var machine = _.clone($scope.machine);

    // Remove currently unused properties
    delete machine.UnavailMsg;
    delete machine.UnavailTill;

    machine.Price = parseFloat(machine.Price);
    if (machine.ReservationPriceStart) {
      machine.ReservationPriceStart = parseFloat(machine.ReservationPriceStart);
    } else {
      machine.ReservationPriceStart = null;
    }
    if (machine.ReservationPriceHourly) {
      machine.ReservationPriceHourly = parseFloat(machine.ReservationPriceHourly);
    } else {
      machine.ReservationPriceHourly = null;
    }

    if (!machine.TypeId) {
      machine.TypeId = 0;
    }
    machine.TypeId = parseInt(machine.TypeId);

    $http({
      method: 'PUT',
      url: '/api/machines/' + $scope.machine.Id,
      headers: {'Content-Type': 'application/json' },
      data: machine,
      transformRequest: function(data) {
        console.log('Machine data to send:', data);
        return JSON.stringify(data);
      },
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      toastr.success('Update successful');
    })
    .error(function(message, statusCode) {
      if (statusCode === 400 && message.indexOf('Dimensions') >= 0) {
        toastr.error(message);
      } else if (statusCode === 400 && message.indexOf('Found machine with same netswitch host') >= 0) {
        toastr.error(message);
      } else {
        toastr.error('Failed to update');
      }
    });
  }; // updateMachine()

  // cf. http://stackoverflow.com/q/17922557/485185
  // There is also a plugin for <input type="file"> on change events.
  $scope.machineImageLoad = function(o) {
    var files = o.files;
    if (files) {
      var f = files[0];
      var reader = new FileReader();
      reader.onloadend = function() {
        $scope.$apply(function() {
          $scope.machineImageNewFile = reader.result;
          $scope.machineImageNewFileName = f.name;
          $scope.machineImageNewFileSize = f.size;
        });
      };
      reader.readAsDataURL(f);
    }
  };

  $scope.machineImageReplace = function() {
    toastr.info('Uploading machine image...');
    $http({
      method: 'POST',
      url: '/api/machines/' + $scope.machine.Id + '/image',
      data: {
        Filename: $scope.machineImageNewFileName,
        Image: $scope.machineImageNewFile
      },
      params: {
        ac: new Date().getTime()
      },
    })
    .success(function(){
      toastr.success('Machine image successfully uploaded');
      $scope.loadMachine();
    })
    .error(function(){
      toastr.error('Uploading machine image failed');
    });
  };

  function configCountdown(seconds, cb) {
    if (seconds >= 0) {
      $timeout(function() {
        $scope.netswitchConfigStatus = 'Updating config... (' + seconds + ' s)';
        configCountdown(seconds - 1, cb);
      }, 1000);
    } else {
      cb();
    }
  }

  $scope.applyConfig = function() {
    if ($scope.unsavedChanges) {
      toastr.error('Please save before continuing.');
      return;
    }
    api.prompt('The power switch upgrade takes around 10 minutes. It enables EASY LAB integration. ', function() {
      $scope.loading = true;
      $http({
        method: 'POST',
        url: '/api/machines/' + $scope.machine.Id + '/apply_config',
        data: {
          location: $cookies.get('location')
        }
      })
      .success(function(){
        $scope.loading = false;
        toastr.success('Updating config...');
        configCountdown(180, function() {
          toastr.success('Configuration pushed.  Switch will be usable in about 5 minutes!');
          $scope.netswitchConfigStatus = undefined;
        });
      })
      .error(function(){
        $scope.loading = false;
        toastr.error('An Error occurred.  Please try again later.');
      });
    });
  };

  var socket;

  function wsConnect() {
    var locationId = $cookies.get('location');
    var t0 = new Date();

    if (socket) {
      socket.onclose = function () {}; // disable onclose handler first
      socket.close();
      socket = undefined;
    }

    //const uid = reactor.evaluateToJS(getters.getUid);
    const host = window.location.host;
    const protocol = host === 'easylab.io' ? 'wss' : 'ws';
    socket = new WebSocket(protocol + '://' + host + '/api/users/' + $scope.user.Id + '/dashboard/ws?location=' + locationId);
    socket.onmessage = function(e) {
      var data = JSON.parse(e.data);
      console.log(data);
      if (data.UserMessage && data.UserMessage.Error) {
        toastr.error(data.UserMessage.Error);
      }
      if (data.UserMessage && data.UserMessage.Info) {
        toastr.info(data.UserMessage.Info);
      }
      if (data.UserMessage && data.UserMessage.Warning) {
        toastr.warn(data.UserMessage.Warning);
      }
    };
    socket.onclose = function(e) {
      console.log('WebSocket closed, e=', e);
      var duration = new Date() - t0;
      socket = null;
      console.log('reconnecting in 5 s...');
      window.setTimeout(function() {
        wsConnect();
      }, 5000);
    };
    socket.onerror = function(e) {
      console.log('websocket error:', e);
    };
  }

  $http({
    method: 'GET',
    url: '/api/users/current',
    params: {
      ac: new Date().getTime()
    }
  })
  .success(function(user) {
    $scope.user = user;
    console.log('$scope.user=', $scope.user);
    wsConnect();
  })
  .error(function() {
    toastr.error('Error obtaining user.  Please refresh.');
  });

}]); // app.controller

})(); // closure
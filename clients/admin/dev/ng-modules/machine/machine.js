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
 ['$scope', '$routeParams', '$http', '$location', '$filter', 'randomToken', 
 function($scope, $routeParams, $http, $location, $filter, randomToken) {

  $scope.machine = {
    Id: $routeParams.machineId
  };

  $scope.machineImageFile = undefined;
  $scope.machineImageNewFile = undefined;
  $scope.machineImageNewFileName = undefined;
  $scope.machineImageNewFileSize = undefined;

  $scope.loadLocations = function() {
    $http({
      method: 'GET',
      url: '/api/locations'
    })
    .success(function(data) {
      $scope.locations = data;
    })
    .error(function() {
      toastr.error('Failed to get locations');
    });
  };


  $scope.loadMachine = function(machineId) {
    $http({
      method: 'GET',
      url: '/api/machines/' + machineId,
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      $scope.machine = data;
      $scope.machine.Price = $filter('currency')($scope.machine.Price,'',2);
      $scope.getNetSwitchMapping();
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

  $scope.loadConnectedMachines = function(machineId) {
    $http({
      method: 'GET',
      url: '/api/machines/' + machineId + '/connections',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(machineList) {
      console.log(machineList);
      $scope.connectedMachines = machineList.Data;
    })
    .error(function() {
      toastr.error('Failed to load connected machines');
    });
  };

  $scope.loadConnectableMachines = function(machineId) {
    $http({
      method: 'GET',
      url: '/api/machines/' + machineId + '/connectable',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(machineList) {
      $scope.connectableMachines = machineList.Data;
    })
    .error(function() {
      toastr.error('Failed to load connectable machines');
    });
  };

  $scope.loadLocations();
  $scope.loadMachineTypes();
  $scope.loadMachine($scope.machine.Id);
  $scope.loadConnectedMachines($scope.machine.Id);
  $scope.loadConnectableMachines($scope.machine.Id);

  $scope.updateMachine = function() {

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

    if (!machine.LocationId) {
      machine.LocationId = 0;
    }
    machine.LocationId = parseInt(machine.LocationId);

    if (!machine.TypeId) {
      machine.TypeId = 0;
    }
    machine.TypeId = parseInt(machine.TypeId);

    if (!machine.LocationId) {
      toastr.error('Please specify a location.');
      return;
    }

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
      } else {
        toastr.error('Failed to update');
      }
    });
  }; // updateMachine()

  $scope.deleteMachinePrompt = function() {
    var token = randomToken.generate();
    vex.dialog.prompt({
      message: 'Enter <span class="delete-prompt-token">' + 
       token + '</span> to delete',
      placeholder: 'Token',
      callback: $scope.deleteMachinePromptCallback.bind(this, token)
    });
  };

  $scope.deleteMachinePromptCallback = function(expectedToken, value) {
    if (value) {    
      if (value === expectedToken) {
        $scope.deleteMachine();
      } else {
        toastr.error('Wrong token');
      }
    } else if (value !== false) {
      toastr.error('No token');
    }
  };

  $scope.deleteMachine = function() {
    $http({
      method: 'DELETE',
      url: '/api/machines/' + $scope.machine.Id,
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function() {
      toastr.success('Machine deleted');
      $location.path('/machines');
    })
    .error(function() {
      toastr.error('Failed to delete machine');
    });
  };

  // The NetSwitch
  // TODO: Create another module for this

  $scope.getNetSwitchMapping = function() {
    $http({
      method: 'GET',
      url: '/api/netswitch/' + $scope.machine.Id,
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(mappingModel) {
      $scope.netSwitchMapping = mappingModel;
      $scope.netSwitchMappingChanged = false;
    }); // no error - the mapping will just not be visible
  };

  $scope.createNetSwitchMapping = function() {
    $http({
      method: 'POST',
      url: '/api/netswitch',
      params: {
        mid: $scope.machine.Id,
        ac: new Date().getTime()
      }
    })
    .success(function(mappingId) {
      $scope.getNetSwitchMapping();
    })
    .error(function() {
      toastr.error('Failed to create NetSwitch mapping');
    });
  };

  $scope.deleteNetSwitchMappingPrompt = function() {
    var token = randomToken.generate();
    vex.dialog.prompt({
      message: 'Enter <span class="delete-prompt-token">' + 
       token + '</span> to delete',
      placeholder: 'Token',
      callback: $scope.deleteNetSwitchMappingPromptCallback.bind(this, token)
    });
  };

  $scope.deleteNetSwitchMappingPromptCallback = function(expectedToken, value) {
    if (value) {    
      if (value === expectedToken) {
        $scope.deleteNetSwitchMapping();
      } else {
        toastr.error('Wrong token');
      }
    } else if (value !== false) {
      toastr.error('No token');
    }
  };

  $scope.deleteNetSwitchMapping = function() {
    $http({
      method: 'DELETE',
      url: '/api/netswitch/' + $scope.machine.Id,
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function() {
      toastr.success('Mapping deleted');
      delete $scope.netSwitchMapping;
    })
    .error(function() {
      toastr.error('Failed to delete mapping');
    });
  };

  // Update the mapping with fresh IP
  $scope.updateNetSwitchMapping = function () {
    if ($scope.netSwitchMapping) {
      $http({
        method: 'PUT',
        url: '/api/netswitch/' + $scope.machine.Id,
        headers: {'Content-Type': 'application/json' },
        data: $scope.netSwitchMapping,
        transformRequest: function(data) {
          return JSON.stringify(data);
        },
        params: {
          ac: new Date().getTime()
        }
      })
      .success(function() {
        $scope.netSwitchMappingChanged = false;
        toastr.success('NetSwitch mapping updated');
      })
      .error(function() {
        toastr.error('Failed to update NetSwitch mapping');
      });
    }
  };

  // Connected machine stuff
  // TODO: Put this in separate module / file
  $scope.addConnectedMachine = function() {
    var connMachineId = $('#connectable-machine-select').val();
    if (!connMachineId) {
      toastr.error('Please select a machine to connect');
      return;
    }

    // Store connectable
    var connMachine = {
      Id: $('#connectable-machine-select').val(),
      Name: $('#connectable-machine-select option:selected').text()
    };

    // Remove the connectable option so we do not repeat ourselves
    for (var i = 0; i < $scope.connectableMachines.length; i++) {
      if (parseInt($scope.connectableMachines[i].Id) === parseInt(connMachine.Id)) {
        $scope.connectableMachines.splice(i, 1);
        break;
      }
    }

    // Add machine to the connected machine list
    if (!$scope.connectedMachines) {
      $scope.connectedMachines = [];
    }
    $scope.connectedMachines.push(connMachine);

    // And also update the machine.ConnectedMachines string based array
    var str = '';
    for (i = 0; i < $scope.connectedMachines.length; i++) {
      str += $scope.connectedMachines[i].Id + ',';
    }
    str = '[' + str.substr(0, str.length - 1) + ']';
    $scope.machine.ConnectedMachines = str;

    $scope.updateMachine();
  };

  $scope.removeConnectedMachinePrompt = function(connMachineId) {
    var token = randomToken.generate();
    vex.dialog.prompt({
      message: 'Enter <span class="delete-prompt-token">' + 
       token + '</span> to remove',
      placeholder: 'Token',
      callback: $scope.removeConnectedMachinePromptCallback.
        bind(this, token, connMachineId)
    });
  };

  $scope.removeConnectedMachinePromptCallback = 
    function(expectedToken, connMachineId, value) {
    if (value) {    
      if (value === expectedToken) {
        $scope.removeConnectedMachine(connMachineId);
      } else {
        toastr.error('Wrong token');
      }
    } else if (value !== false) {
      toastr.error('No token');
    }
  };

  $scope.removeConnectedMachine = function(connMachineId) {

    // Search for connected machine with the ID so we can swap
    // move it from connected machines to connectable machines
    for (var i = 0; i < $scope.connectedMachines.length; i++) {
      if (parseInt($scope.connectedMachines[i].Id) === parseInt(connMachineId)) {
        if (!$scope.connectableMachines) {
          $scope.connectableMachines = [];
        }
        $scope.connectableMachines.push($scope.connectedMachines[i]);
        $scope.connectedMachines.splice(i, 1);
        $scope.$apply();
        break;
      }
    }

    // And also update the machine.ConnectedMachines string based array
    // User will have to press `Save` to update the database
    var str = '';
    for (i = 0; i < $scope.connectedMachines.length; i++) {
      str += $scope.connectedMachines[i].Id + ',';
    }
    str = '[' + str.substr(0, str.length - 1) + ']';
    $scope.machine.ConnectedMachines = str;

    $scope.updateMachine();
  };

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
    })
    .error(function(){
      toastr.error('Uploading machine image failed');
    });
  };

}]); // app.controller

})(); // closure
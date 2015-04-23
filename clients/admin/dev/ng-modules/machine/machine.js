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

  $scope.machineImageNewFile = undefined;
  $scope.machineImageNewFileName = undefined;
  $scope.machineImageNewFileSize = undefined;

  $http({
    method: 'GET',
    url: '/api/machines/' + $scope.machine.Id,
    params: {
      ac: new Date().getTime()
    }
  })
  .success(function(data) {
    $scope.machine = data;
    $scope.machine.Price = $filter('currency')($scope.machine.Price,'',2);
    $scope.getHexabusMapping();
  })
  .error(function() {
    toastr.error('Failed to get machine');
  });

  $scope.updateMachine = function() {

    // Make a clone of the machine model
    var machine = _.clone($scope.machine);

    // Remove currently unused properties
    delete machine.UnavailMsg;
    delete machine.UnavailTill;
    delete machine.Image;

    machine.Price = parseFloat(machine.Price);

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
      //$scope.updateHexabusMapping();
    })
    .error(function(data) {
      console.log(data);
      toastr.error('Failed to update');
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

  // TODO: move hexabus mapping to another module

  $scope.getHexabusMapping = function() {
    $http({
      method: 'GET',
      url: '/api/hexabus/' + $scope.machine.Id,
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(mappingModel) {
      $scope.hexabusMapping = mappingModel;
      $scope.hexabusMappingChanged = false;
    }); // no error - the mapping will just not be visible
  };

  // Create a field for entering the IP
  $scope.createHexabusMapping = function() {
    $http({
      method: 'POST',
      url: '/api/hexabus',
      params: {
        mid: $scope.machine.Id,
        ac: new Date().getTime()
      }
    })
    .success(function(mappingId) {
      $scope.getHexabusMapping();
    })
    .error(function() {
      toastr.error('Failed to create hexabus mapping');
    });
    
  };

  $scope.deleteHexabusMappingPrompt = function() {
    var token = randomToken.generate();
    vex.dialog.prompt({
      message: 'Enter <span class="delete-prompt-token">' + 
       token + '</span> to delete',
      placeholder: 'Token',
      callback: $scope.deleteHexabusMappingPromptCallback.bind(this, token)
    });
  };

  $scope.deleteHexabusMappingPromptCallback = function(expectedToken, value) {
    if (value) {    
      if (value === expectedToken) {
        $scope.deleteHexabusMapping();
      } else {
        toastr.error('Wrong token');
      }
    } else if (value !== false) {
      toastr.error('No token');
    }
  };

  $scope.deleteHexabusMapping = function() {
    $http({
      method: 'DELETE',
      url: '/api/hexabus/' + $scope.machine.Id,
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function() {
      toastr.success('Mapping deleted');
      delete $scope.hexabusMapping;
    })
    .error(function() {
      toastr.error('Failed to delete mapping');
    });
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
    toastr.info('machineImageReplace()');
    $http({
      method: 'POST',
      url: '/api/machines/' + $scope.machine.Id + '/image',
      data: {
        Filename: $scope.machineImageNewFileName,
        Image: $scope.machineImageNewFile
      }
    })
  };

  // Update the mapping with fresh IP
  $scope.updateHexabusMapping = function () {
    if ($scope.hexabusMapping) {
      $http({
        method: 'PUT',
        url: '/api/hexabus/' + $scope.machine.Id,
        headers: {'Content-Type': 'application/json' },
        data: $scope.hexabusMapping,
        transformRequest: function(data) {
          return JSON.stringify(data);
        },
        params: {
          ac: new Date().getTime()
        }
      })
      .success(function() {
        $scope.hexabusMappingChanged = false;
        toastr.success('Hexabus mapping updated');
      })
      .error(function() {
        toastr.error('Failed to update hexabus mapping');
      });
    }
  };

}]); // app.controller

})(); // closure
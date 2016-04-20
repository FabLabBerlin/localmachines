(function(){

'use strict';

var app = angular.module('fabsmith.admin.machines', ['ngRoute', 'ngCookies', 'fabsmith.admin.api']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/machines', {
    templateUrl: 'ng-modules/machines/machines.html',
    controller: 'MachinesCtrl'
  });
}]); // app.config

app.controller('MachinesCtrl',
 ['$scope', '$http', '$location', '$cookies', 'api',
 function($scope, $http, $location, $cookies, api) {

  $scope.machines = [];

  api.loadMachines(function(resp) {
    $scope.machines = resp.machines;
  });

  $scope.addMachinePrompt = function() {
    vex.dialog.prompt({
      message: 'Enter machine name',
      placeholder: 'Machine name',
      callback: $scope.machinePromptCallback
    });
  };

  $scope.machinePromptCallback = function(value) {
    if (value) {    
      $scope.addMachine(value);
    } else if (value !== false) {
      toastr.error('No machine name');
    }
  };

  $scope.addMachine = function(machineName) {
    $http({
      method: 'POST',
      url: '/api/machines',
      params: {
        location: $cookies.locationId,
        mname: machineName
      }
    })
    .success(function(data) {
      $scope.editMachine(data.MachineId);
    })
    .error(function() {
      toastr.error('Failed to create machine');
    });
  };

  $scope.editMachine = function(id) {
    $location.path('/machine/' + id);
  };

  $scope.setShowArchived = function(show) {
    $scope.showArchived = show;
  };

}]); // app.controller

app.filter('machinesFilter', function() {
  return function(machines, scope) {
    return _.filter(machines, function(machine) {
      return scope.showArchived || !machine.Archived;
    });
  }
});

})(); // closure
(function(){

'use strict';

var app = angular.module('fabsmith.admin.machines', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/machines', {
    templateUrl: 'ng-modules/machines/machines.html',
    controller: 'MachinesCtrl'
  });
}]); // app.config

app.controller('MachinesCtrl', ['$scope', '$http', '$location', '$cookieStore', 
 function($scope, $http, $location, $cookieStore) {

  $scope.machines = [];

  $http({
    method: 'GET',
    url: '/api/machines',
    params: {
      ac: new Date().getTime()
    }
  })
  .success(function(data) {
    $scope.machines = data;
  })
  .error(function() {
    toastr.error('Failed to get machines');
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
        mname: machineName,
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      $scope.editMachine(data.MachineId);
    })
    .error(function() {
      toastr.error('Failed to create machine');
    });
  };

  $scope.editMachine = function(machineId) {
    $location.path('/machine/' + machineId);
    //$scope.$apply();
  };

}]); // app.controller

})(); // closure
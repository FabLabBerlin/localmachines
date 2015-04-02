(function(){

'use strict';

var app = angular.module('fabsmith.admin.machines', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/machines', {
    templateUrl: 'ng-modules/machines/machines.html',
    controller: 'MachinesCtrl'
  });
}]); // app.config

app.controller('MachinesCtrl', ['$scope', '$http', '$location', 
 function($scope, $http, $location) {

  $scope.machines = [
    {
      Id: 1, 
      Name: 'Laser Cutter',
      Price: 10,
      PriceUnit: 'minute'
    }, {
      Id: 2, 
      Name: '3D Priter',
      Price: 10,
      PriceUnit: 'minute'
    }, {
      Id: 3, 
      Name: 'CNC Mill',
      Price: 10,
      PriceUnit: 'hour'
    }
  ];

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
        anticache: new Date().getTime()
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
    $scope.$apply();
  };

}]); // app.controller

})(); // closure
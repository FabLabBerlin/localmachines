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

  $scope.addMachinePrompt = function() {
    vex.dialog.prompt({
      message: 'Enter machine name',
      placeholder: 'Machine name',
      callback: $scope.machinePromptCallback
    });
  };

  $scope.machinePromptCallback = function(value) {
    if (value) {
      toastr.success('Creating machine');
    } else if (value !== false) {
      toastr.error('No machine name');
    }
  };

}]); // app.controller

})(); // closure
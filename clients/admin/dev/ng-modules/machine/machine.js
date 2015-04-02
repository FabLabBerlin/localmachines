(function(){

'use strict';

var app = angular.module('fabsmith.admin.machine', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/machine/:machineId', {
    templateUrl: 'ng-modules/machine/machine.html',
    controller: 'MachineCtrl'
  });
}]); // app.config

app.controller('MachineCtrl', ['$scope', '$routeParams', '$http', '$location', 
 function($scope, $routeParams, $http, $location) {

  $scope.machine = {
    Id: $routeParams.machineId
  };

  console.log('machine.Id: ' + $scope.machine.Id);

}]); // app.controller

})(); // closure
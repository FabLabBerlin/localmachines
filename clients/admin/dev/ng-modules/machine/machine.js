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
    Id: $routeParams.machineId, 
    Name: 'Test Name',
    Shortname: 'TN',
    Price: 10,
    PriceUnit: 'hour',
    Description: 'Lorem ipsum scriptum desc',
    ImageSrc: 'assets/img/img-3d-printer.svg',
    ImageName: 'img-3d-printer.svg',
    ImageSize: '1.9K'
  };

  console.log('machine.Id: ' + $scope.machine.Id);

}]); // app.controller

})(); // closure
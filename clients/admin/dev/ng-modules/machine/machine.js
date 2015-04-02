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

  $http({
    method: 'GET',
    url: '/api/machines/' + $scope.machine.Id,
    params: {
      anticache: new Date().getTime()
    }
  })
  .success(function(data) {
    console.log(data);
    $scope.machine = data;
  })
  .error(function() {
    toastr.error('Failed to get machine');
  });

}]); // app.controller

})(); // closure
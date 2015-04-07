(function(){

'use strict';

var app = angular.module('fabsmith.admin.machine', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/machine/:machineId', {
    templateUrl: 'ng-modules/machine/machine.html',
    controller: 'MachineCtrl'
  });
}]); // app.config

app.controller('MachineCtrl', 
 ['$scope', '$routeParams', '$http', '$location', '$filter', 
 function($scope, $routeParams, $http, $location, $filter) {

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
    $scope.machine.Price = $filter('currency')($scope.machine.Price,'',2);
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

    console.log(machine);
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
        anticache: new Date().getTime()
      }
    })
    .success(function(data) {
      console.log(data);
      toastr.success('Update successful');
    })
    .error(function(data) {
      console.log(data);
      toastr.error('Failed to update');
    });
  };

}]); // app.controller

})(); // closure
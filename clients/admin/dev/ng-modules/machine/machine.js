(function(){

'use strict';

var app = angular.module('fabsmith.admin.machine', 
 ['ngRoute', 'ngCookies']);

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
  }; // updateMachine()

  $scope.deleteMachinePrompt = function() {
    console.log('delete machine prompt');

    // You have to add the <random-token></random-token> directive somewhere
    // in HTML in order to make this work
    $scope.generateRandomToken();

    vex.dialog.prompt({
      // Unfortunately it is not possible to parse directives inside
      // vex messages, so we just get the random token
      message: 'Enter <span class="delete-prompt-token">' + 
       $scope.randomToken + '</span> to delete',
      placeholder: 'Token',
      callback: $scope.deleteMachinePromptCallback
    });
  };

  $scope.deleteMachinePromptCallback = function(value) {
    if (value) {    
      if (value === $scope.randomToken) {
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
        anticache: new Date().getTime()
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

}]); // app.controller

})(); // closure
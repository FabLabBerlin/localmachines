(function(){

'use strict';
var app = angular.module('fabsmith.admin.spaces', ['ngRoute', 'ngCookies', 'fabsmith.admin.randomtoken']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/spaces', {
    templateUrl: 'ng-modules/spaces/spaces.html',
    controller: 'SpacesCtrl'
  });
}]); // app.config

app.controller('SpacesCtrl',
 ['$scope', '$routeParams', '$http', '$location', 'randomToken',
 function($scope, $routeParams, $http, $location, randomToken) {

  $scope.spaces = [];

  /*
   *
   * Spaces functions
   *
   */

  function loadSpaces() {
    $http({
      method: 'GET',
      url: '/api/spaces',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      $scope.spaces = _.sortBy(data, function(space) {
        return space.Product.Name;
      });
    })
    .error(function() {
      toastr.error('Failed to get spaces');
    });
  }

  $scope.addSpacePrompt = function() {
    vex.dialog.prompt({
      message: 'Enter space name',
      placeholder: 'Space name',
      callback: $scope.spacePromptCallback
    });
  };

  $scope.spacePromptCallback = function(value) {
    if (value) {    
      $scope.addSpace(value);
    } else if (value !== false) {
      toastr.error('No space name');
    }
  };

  $scope.addSpace = function(name) {
    $http({
      method: 'POST',
      url: '/api/spaces',
      params: {
        name: name,
        ac: new Date().getTime()
      }
    })
    .success(function(space) {
      $scope.editSpace(space.Product.Id);
    })
    .error(function() {
      toastr.error('Failed to create space');
    });
  };

  $scope.editSpace = function(id) {
    $location.path('/spaces/' + id);
  };

  /*
   *
   * Space Purchases functions
   *
   */

  function loadSpacePurchases() {
    $http({
      method: 'GET',
      url: '/api/space_purchases',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      $scope.spacePurchases = _.sortBy(data, function(spacePurchase) {
        return spacePurchase.Name;
      });
    })
    .error(function() {
      toastr.error('Failed to get space purchases');
    });
  }

  $scope.addSpacePurchase = function() {
    $http({
      method: 'POST',
      url: '/api/space_purchases',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(spacePurchase) {
      $scope.editSpace(spacePurchase.Id);
    })
    .error(function() {
      toastr.error('Failed to create space purchase');
    });
  };

  $scope.editSpacePurchase = function(id) {
    $location.path('/space_purchases/' + id);
  };

  loadSpaces();
  loadSpacePurchases();

}]); // app.controller

})(); // closure

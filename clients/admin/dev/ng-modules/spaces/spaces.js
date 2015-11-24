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
  $scope.spacesById = {};
  $scope.usersById = {};

  /*
   *
   * Spaces functions
   *
   */

  function loadSpaces() {
    $http({
      method: 'GET',
      url: '/api/products',
      params: {
        ac: new Date().getTime(),
        type: 'space'
      }
    })
    .success(function(data) {
      $scope.spaces = _.sortBy(data, function(space) {
        return space.Product.Name;
      });
      _.each($scope.spaces, function(space) {
        $scope.spacesById[space.Product.Id] = space;
      });
      loadUsers();
    })
    .error(function() {
      toastr.error('Failed to get spaces');
    });
  }

  function loadUsers() {
    $http({
      method: 'GET',
      url: '/api/users',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      $scope.users = _.sortBy(data, function(user) {
        return user.FirstName + ' ' + user.LastName;
      });
      _.each($scope.users, function(user) {
        $scope.usersById[user.Id] = user;
      });
      loadSpacePurchases();
    })
    .error(function() {
      toastr.error('Failed to get reservations');
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
      url: '/api/products',
      params: {
        name: name,
        ac: new Date().getTime(),
        type: 'space'
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
      $scope.spacePurchases = _.map($scope.spacePurchases, function(sp) {
        var space = $scope.spacesById[sp.ProductId];
        if (space) {
          sp.Product = space.Product;
        }
        var user = $scope.usersById[sp.UserId];
        sp.User = user;
        sp.TimeStartLocal = moment(sp.TimeStart).tz('Europe/Berlin').format('YYYY-MM-DD HH:mm');
        sp.TimeEndLocal = moment(sp.TimeEnd).tz('Europe/Berlin').format('YYYY-MM-DD HH:mm');
        return sp;
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
      $scope.editSpacePurchase(spacePurchase.Id);
    })
    .error(function() {
      toastr.error('Failed to create space purchase');
    });
  };

  $scope.editSpacePurchase = function(id) {
    $location.path('/space_purchases/' + id);
  };

  loadSpaces();

}]); // app.controller

})(); // closure

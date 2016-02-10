(function(){

'use strict';
var app = angular.module('fabsmith.admin.spaces', ['ngRoute', 'ngCookies', 'fabsmith.admin.randomtoken', 'fabsmith.admin.api']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/spaces', {
    templateUrl: 'ng-modules/spaces/spaces.html',
    controller: 'SpacesCtrl'
  });
}]); // app.config

app.controller('SpacesCtrl',
 ['$scope', '$http', '$location', '$cookies', 'randomToken', 'api',
 function($scope, $http, $location, $cookies, randomToken, api) {

  $scope.spaces = [];
  $scope.spacesById = {};
  $scope.usersById = {};

  $scope.addSpacesProductPromptCallback = function(value) {
    if (value) {    
      $scope.addSpacesProduct(value);
    } else if (value !== false) {
      toastr.error('No product name');
    }
  };

  $scope.addSpacesProductPrompt = function() {
    vex.dialog.prompt({
      message: 'Enter spaces product name',
      placeholder: 'Product name',
      callback: $scope.addSpacesProductPromptCallback
    });
  };

  $scope.addSpacesProduct = function(name) {
    $http({
      method: 'POST',
      url: '/api/products',
      params: {
        location: $cookies.locationId,
        name: name,
        ac: new Date().getTime(),
        type: 'space'
      }
    })
    .success(function(data) {
      $scope.editSpacesProduct(data.Product.Id);
    })
    .error(function() {
      toastr.error('Failed to create product');
    });
  };

  $scope.editSpacesProduct = function(id) {
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
      url: '/api/purchases',
      params: {
        ac: new Date().getTime(),
        type: 'space'
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
      url: '/api/purchases',
      params: {
        ac: new Date().getTime(),
        type: 'space'
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

  api.loadSpaces(function(spacesData) {
    $scope.spaces = spacesData.spaces;
    console.log('$scope.spaces= ', $scope.spaces);
    $scope.spacesById = spacesData.spacesById;
    api.loadUsers(function(userData) {
      $scope.users = userData.users;
      $scope.usersById = userData.usersById;
      loadSpacePurchases();
    });
  });

}]); // app.controller

})(); // closure

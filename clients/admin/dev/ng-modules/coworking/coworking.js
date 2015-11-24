(function(){

'use strict';
var app = angular.module('fabsmith.admin.coworking', ['ngRoute', 'ngCookies', 'fabsmith.admin.randomtoken']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/coworking', {
    templateUrl: 'ng-modules/coworking/coworking.html',
    controller: 'CoworkingCtrl'
  });
}]); // app.config

app.controller('CoworkingCtrl',
 ['$scope', '$routeParams', '$http', '$location', 'randomToken',
 function($scope, $routeParams, $http, $location, randomToken) {

  $scope.tables = [];
  $scope.tablesById = {};
  $scope.usersById = {};

  /*
   *
   * Tables functions
   *
   */

  function loadTables() {
    $http({
      method: 'GET',
      url: '/api/products',
      params: {
        ac: new Date().getTime(),
        type: 'co-working'
      }
    })
    .success(function(data) {
      $scope.tables = _.sortBy(data, function(table) {
        return table.Product.Name;
      });
      _.each($scope.tables, function(table) {
        $scope.tablesById[table.Product.Id] = table;
      });
      loadUsers();
    })
    .error(function() {
      toastr.error('Failed to get tables');
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
      loadRentals();
    })
    .error(function() {
      toastr.error('Failed to get users');
    });
  }

  /*
   *
   * Table Purchases functions
   *
   */

  function loadRentals() {
    $http({
      method: 'GET',
      url: '/api/purchases',
      params: {
        ac: new Date().getTime(),
        type: 'co-working'
      }
    })
    .success(function(data) {
      $scope.rentals = _.sortBy(data, function(rental) {
        return rental.Name;
      });
      $scope.rentals = _.map($scope.rentals, function(sp) {
        var table = $scope.tablesById[sp.ProductId];
        if (table) {
          sp.Product = table.Product;
        }
        var user = $scope.usersById[sp.UserId];
        sp.User = user;
        sp.TimeStartLocal = moment(sp.TimeStart).tz('Europe/Berlin').format('YYYY-MM-DD HH:mm');
        sp.TimeEndLocal = moment(sp.TimeEnd).tz('Europe/Berlin').format('YYYY-MM-DD HH:mm');
        return sp;
      });
    })
    .error(function() {
      toastr.error('Failed to get table purchases');
    });
  }

  $scope.addRental = function() {
    $http({
      method: 'POST',
      url: '/api/purchases',
      params: {
        ac: new Date().getTime(),
        type: 'co-working'
      }
    })
    .success(function(rental) {
      $scope.editRental(rental.Id);
    })
    .error(function() {
      toastr.error('Failed to create table purchase');
    });
  };

  $scope.editRental = function(id) {
    $location.path('/coworking/rentals/' + id);
  };

  loadTables();

}]); // app.controller

})(); // closure

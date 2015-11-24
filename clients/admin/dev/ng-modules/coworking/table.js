(function(){

'use strict';
var app = angular.module('fabsmith.admin.coworking.table', ['ngRoute', 'ngCookies', 'fabsmith.admin.randomtoken']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/coworking/tables/:id', {
    templateUrl: 'ng-modules/coworking/table.html',
    controller: 'TableCtrl'
  });
}]); // app.config

app.controller('TableCtrl',
 ['$scope', '$routeParams', '$http', '$location', 'randomToken',
 function($scope, $routeParams, $http, $location, randomToken) {

  $scope.table = {
    Product: {
      Id: $routeParams.id
    }
  };

  function loadTable() {
    $http({
      method: 'GET',
      url: '/api/products/' + $scope.table.Product.Id,
      params: {
        ac: new Date().getTime(),
        type: 'co-working'
      }
    })
    .success(function(table) {
      $scope.table = table;
    })
    .error(function(data, status) {
      toastr.error('Failed to load tables data');
    });
  }

  $scope.updateTable = function() {
    $http({
      method: 'PUT',
      url: '/api/products/' + $scope.table.Product.Id + '?type=co-working',
      headers: {'Content-Type': 'application/json' },
      data: $scope.table,
      transformRequest: function(data) {
        var transformed = {
          Product: _.extend({}, data.Product)
        };
        transformed.Product.Id = parseInt(data.Product.Id);
        transformed.Product.Price = parseFloat(data.Product.Price);
        return JSON.stringify(transformed);
      },
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      toastr.success('Update successful');
    })
    .error(function(data) {
      console.log(data);
      toastr.error('Failed to update');
    });
  };

  loadTable();

}]); // app.controller

})(); // closure

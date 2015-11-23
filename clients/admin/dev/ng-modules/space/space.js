(function(){

'use strict';
var app = angular.module('fabsmith.admin.space', ['ngRoute', 'ngCookies', 'fabsmith.admin.randomtoken']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/spaces/:id', {
    templateUrl: 'ng-modules/space/space.html',
    controller: 'SpaceCtrl'
  });
}]); // app.config

app.controller('SpaceCtrl',
 ['$scope', '$routeParams', '$http', '$location', 'randomToken',
 function($scope, $routeParams, $http, $location, randomToken) {

  $scope.space = {
    Product: {
      Id: $routeParams.id
    }
  };

  function loadSpace() {
    $http({
      method: 'GET',
      url: '/api/spaces/' + $scope.space.Product.Id,
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(space) {
      $scope.space = space;
    })
    .error(function(data, status) {
      toastr.error('Failed to load user data');
    });
  };


  $scope.updateSpace = function() {
    $http({
      method: 'PUT',
      url: '/api/spaces/' + $scope.space.Product.Id,
      headers: {'Content-Type': 'application/json' },
      data: $scope.space,
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

  loadSpace();

}]); // app.controller

})(); // closure

(function(){

'use strict';
var app = angular.module('fabsmith.admin.coworking.product', ['ngRoute', 'ngCookies', 'fabsmith.admin.randomtoken']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/coworking/products/:id', {
    templateUrl: 'ng-modules/coworking/coworking-product.html',
    controller: 'CoworkingProductCtrl'
  });
}]); // app.config

app.controller('CoworkingProductCtrl',
 ['$scope', '$routeParams', '$http', '$location', 'randomToken',
 function($scope, $routeParams, $http, $location, randomToken) {

  $scope.product = {
    Product: {
      Id: $routeParams.id
    }
  };

  function loadProduct() {
    $http({
      method: 'GET',
      url: '/api/products/' + $scope.product.Product.Id,
      params: {
        ac: new Date().getTime(),
        type: 'co-working'
      }
    })
    .success(function(product) {
      $scope.product = product;
      $scope.product.Product.PriceUnit = 'month';
    })
    .error(function() {
      toastr.error('Failed to load product data');
    });
  }

  $scope.updateProduct = function() {
    $http({
      method: 'PUT',
      url: '/api/products/' + $scope.product.Product.Id + '?type=co-working',
      headers: {'Content-Type': 'application/json' },
      data: $scope.product,
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
      $location.path('/coworking');
    })
    .error(function(data) {
      console.log(data);
      toastr.error('Failed to update');
    });
  };

  loadProduct();

}]); // app.controller

})(); // closure

(function(){

'use strict';

var app = angular.module("fabsmith.admin.productlist", []);

app.directive('productlist', function() {
  return {
    restrict: 'E',
    templateUrl: 'ng-modules/productlist/productlist.html',
    controller: function($scope, $element, $attrs, $http, $location) {
      $scope.addProductPrompt = function() {
        vex.dialog.prompt({
          message: 'Enter product name',
          placeholder: 'Product name',
          callback: $scope.addProductPromptCallback
        });
      };

      $scope.addProductPromptCallback = function(value) {
        if (value) {    
          $scope.addProduct(value);
        } else if (value !== false) {
          toastr.error('No product name');
        }
      };

      $scope.addProduct = function(name) {
        $http({
          method: 'POST',
          url: '/api/products',
          params: {
            name: name,
            ac: new Date().getTime(),
            type: $attrs.type
          }
        })
        .success(function(data) {
          $scope.editProduct(data.Product.Id);
        })
        .error(function() {
          toastr.error('Failed to create product');
        });
      };

      $scope.editProduct = function(id) {
        switch ($attrs.type) {
        case 'co-working':
          $location.path('/coworking/tables/' + id);
          break;
        case 'space':
          $location.path('/spaces/' + id);
          break;
        default:
          console.log('Product list: $attrs.type = ', $attrs.type);
          toastr.error('Product list: unknown product type');
        }
      };

      $http({
        method: 'GET',
        url: '/api/products',
        params: {
          ac: new Date().getTime(),
          type: $attrs.type
        }
      })
      .success(function(products) {
        $scope.products = products;
        switch ($attrs.type) {
        case 'co-working':
          $scope.products = _.pluck($scope.products, 'Product');
          console.log('$scope.products=', $scope.products);
          break;
        case 'space':
          $scope.products = _.pluck($scope.products, 'Product');
          break;
        default:
          console.log('Product list: $attrs.type = "' + $attrs.type + '"');
          toastr.error('Product list: unknown product type');
        }
      })
      .error(function() {
        toastr.error('Failed to get products');
      });
    },
    link: function($scope, $element, $attrs) {
    }
  };
});

})(); // closure

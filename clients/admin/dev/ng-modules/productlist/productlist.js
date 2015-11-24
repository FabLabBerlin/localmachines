(function(){

'use strict';

angular.module("fabsmith.admin.productlist", [])
.directive('productlist', function() {
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
          toastr.error('Failed to create space');
        });
      };

      $scope.editProduct = function(id) {
        switch ($attrs.type) {
        case 'space':
          $location.path('/spaces/' + id);
          break;
        default:
          console.log('Product list: $attrs.type = ', $attrs.type);
          toastr.error('Product list: unknown product type');
        }
        
      }
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
        if ($attrs.type === 'space') {
          $scope.products = _.pluck($scope.products, 'Product');
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

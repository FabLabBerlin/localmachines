(function(){

'use strict';
var app = angular.module('fabsmith.admin.coworking', 
  ['ngRoute', 'ngCookies', 'fabsmith.admin.randomtoken']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/coworking', {
    templateUrl: 'ng-modules/coworking/coworking.html',
    controller: 'CoworkingCtrl'
  });
}]); // app.config

app.controller('CoworkingCtrl',
 ['$scope', '$routeParams', '$http', '$location', '$cookies', 'randomToken',
 function($scope, $routeParams, $http, $location, $cookies, randomToken) {

  $scope.products = [];
  $scope.productsById = {};
  $scope.usersById = {};

  $scope.loadCoworkingProducts = function(success, error) {
    $http({
      method: 'GET',
      url: '/api/products',
      params: {
        location: $cookies.location,
        ac: new Date().getTime(),
        type: 'co-working'
      }
    })
    .success(function(products) {
      if (success) {
        success(products);
      }
    })
    .error(function() {
      if (error) {
        error();
      }
    });
  };

  $scope.addCoworkingProductPrompt = function() {
    vex.dialog.prompt({
      message: 'Enter coworking product name',
      placeholder: 'Product name',
      callback: $scope.addCoworkingProductPromptCallback
    });
  };

  $scope.addCoworkingProductPromptCallback = function(value) {
    if (value) {    
      $scope.addCoworkingProduct(value);
    } else if (value !== false) {
      toastr.error('No product name');
    }
  };

  $scope.addCoworkingProduct = function(name) {
    $http({
      method: 'POST',
      url: '/api/products',
      params: {
        location: $cookies.location,
        name: name,
        ac: new Date().getTime(),
        type: 'co-working'
      }
    })
    .success(function(data) {
      $scope.editCoworkingProduct(data.Product.Id);
    })
    .error(function() {
      toastr.error('Failed to create product');
    });
  };

  $scope.editCoworkingProduct = function(id) {
    $location.path('/coworking/products/' + id);
  };

  $scope.loadUsers = function(success, error) {
    $http({
      method: 'GET',
      url: '/api/users',
      params: {
        ac: new Date().getTime(),
        location: $cookies.location
      }
    })
    .success(function(users) {
      if (success) {
        success(users);
      }
    })
    .error(function() {
      if (error) {
        error();
      }
    });
  };

  $scope.loadCoworkingPurchases = function(success, error) {
    $http({
      method: 'GET',
      url: '/api/purchases',
      params: {
        location: $cookies.location,
        ac: new Date().getTime(),
        type: 'co-working'
      }
    })
    .success(function(purchases) {
      if (success) {
        success(purchases);
      }
    })
    .error(function() {
      toastr.error('Failed to get coworking purchases');
    });
  };

  $scope.addCoworkingPurchase = function() {
    $http({
      method: 'POST',
      url: '/api/purchases',
      params: {
        location: $cookies.location,
        ac: new Date().getTime(),
        type: 'co-working'
      }
    })
    .success(function(purchase) {
      $scope.editCoworkingPurchase(purchase.Id);
    })
    .error(function() {
      toastr.error('Failed to create coworking purchase');
    });
  };

  $scope.editCoworkingPurchase = function(id) {
    $location.path('/coworking/purchases/' + id);
  };

  $scope.loadCoworkingProducts(function(products) {
    $scope.products = products;
    $scope.products = _.pluck($scope.products, 'Product');
    _.each($scope.products, function(product) {
      $scope.productsById[product.Id] = product;
    });

    // Continue with loading users
    $scope.loadUsers(function(users) {
      $scope.users = _.sortBy(users, function(user) {
        return user.FirstName + ' ' + user.LastName;
      });
      _.each($scope.users, function(user) {
        $scope.usersById[user.Id] = user;
      });

      // Continue with loading coworking purchases
      $scope.loadCoworkingPurchases(function(purchases) {
        $scope.purchases = _.sortBy(purchases, function(purchase) {
          return purchase.Name;
        });
        $scope.purchases = _.map($scope.purchases, function(purchase) {
          var product = $scope.productsById[purchase.ProductId];
          if (product) {
            purchase.Product = product;
          }
          var user = $scope.usersById[purchase.UserId];
          purchase.User = user;
          purchase.TimeStartLocal = moment(purchase.TimeStart).tz('Europe/Berlin').format('YYYY-MM-DD HH:mm');
          purchase.TimeEndLocal = moment(purchase.TimeEnd).tz('Europe/Berlin').format('YYYY-MM-DD HH:mm');
          return purchase;
        });
        console.log($scope.purchases);
      }, function() {
        toastr.error('Failed to load coworking purchases');
      });
    }, function() {
      toastr.error('Failed to get users');
    });
  }, function() {
    toastr.error('Failed to get coworking products');
  });

}]); // app.controller

})(); // closure

(function(){

'use strict';
var app = angular.module('fabsmith.admin.coworking.purchase', 
  ['ngRoute', 'ngCookies', 'fabsmith.admin.randomtoken', 'fabsmith.admin.api']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/coworking/purchases/:id', {
    templateUrl: 'ng-modules/coworking/coworking-purchase.html',
    controller: 'CoworkingPurchaseCtrl'
  });
}]); // app.config

app.controller('CoworkingPurchaseCtrl',
 ['$scope', '$routeParams', '$http', '$location', 'randomToken', 'api', '$cookies',
 function($scope, $routeParams, $http, $location, randomToken, api, $cookies) {

  $scope.purchases = [];
  $scope.purchase = {
    Id: $routeParams.id
  };
  $scope.productsById = {};
  $scope.users = [];
  $scope.usersById = {};

  function loadProducts() {
    $http({
      method: 'GET',
      url: '/api/products',
      params: {
        location: $cookies.get('location'),
        ac: new Date().getTime(),
        type: 'co-working'
      }
    })
    .success(function(data) {
      $scope.products = _.sortBy(data, function(product) {
        return product.Product.Name;
      });
      _.each($scope.products, function(product) {
        $scope.productsById[product.Product.Id] = product;
      });
      loadProductPurchase();
    })
    .error(function() {
      toastr.error('Failed to get products');
    });
  }

  function loadProductPurchase() {
    $http({
      method: 'GET',
      url: '/api/purchases/' + $scope.purchase.Id,
      params: {
        ac: new Date().getTime(),
        type: 'co-working'
      }
    })
    .success(function(sp) { // wtf is sp?
      $scope.purchase = sp;
      var start = moment(sp.TimeStart).tz('Europe/Berlin');
      var end = moment(sp.TimeEnd).tz('Europe/Berlin');
      sp.DateStartLocal = start.format('YYYY-MM-DD');
      sp.DateEndLocal = end.format('YYYY-MM-DD');
      sp.TimeStartLocal = start.format('HH:mm');
      sp.TimeEndLocal = end.format('HH:mm');
      sp.PriceUnit = 'month';
      calculateTotalPrice();
      setTimeout(function() {
        $('.selectpicker').selectpicker('refresh');
      }, 100);
    })
    .error(function(data, status) {
      toastr.error('Failed to load user data');
    });
  }

  function parseInputTimes() {
    var sp = $scope.purchase;
    sp.TimeStart = moment.tz(sp.DateStartLocal + ' ' + sp.TimeStartLocal, 'Europe/Berlin').toDate();
    sp.TimeEnd = moment.tz(sp.DateEndLocal + ' ' + sp.TimeEndLocal, 'Europe/Berlin').toDate();
  }

  function calculateQuantity() {
    console.log('$scope.timeChange()');
    parseInputTimes();
    var start = moment($scope.purchase.TimeStart);
    var end = moment($scope.purchase.TimeEnd);
    var duration = end.unix() - start.unix();
    console.log('duration=', duration);
    var quantity;
    switch ($scope.purchase.PriceUnit) {
    case 'minute':
      quantity = duration / 60;
      break;
    case 'hour':
      quantity = duration / 3600;
      break;
    case 'day':
      quantity = duration / 24 / 3600;
      break;
    case 'month':
      quantity = duration / 24 / 3600 / 12;
      break;
    default:
      return;
    }
    $scope.purchase.Quantity = quantity;
  }
   
  // https://www.artstation.com/artwork/b5zBn
  function calculateTotalPrice() {
    var totalPrice = $scope.purchase.Quantity * $scope.purchase.PricePerUnit;
    $scope.purchase.TotalPrice = totalPrice.toFixed(2);
  }

  $scope.productChange = function() {
    $('#purchase-product').selectpicker('refresh');
    console.log('$scope.productChange()');
    var productId = parseInt($scope.purchase.ProductId);
    var product = $scope.productsById[productId];
    $scope.purchase.PricePerUnit = product.Product.Price;
    $scope.purchase.PriceUnit = product.Product.PriceUnit;
    calculateQuantity();
    calculateTotalPrice();    
  };

  $scope.userChange = function() {
    $('#purchase-user').selectpicker('refresh');
  };

  $scope.timeChange = function() {
    calculateQuantity();
    calculateTotalPrice();
  };

  $scope.priceUnitChange = function() {
    calculateQuantity();
    calculateTotalPrice();
  };

  $scope.save = function() {
    parseInputTimes();
    $http({
      method: 'PUT',
      url: '/api/purchases/' + $scope.purchase.Id + '?type=co-working',
      headers: {'Content-Type': 'application/json' },
      data: $scope.purchase,
      transformRequest: function(data) {
        var transformed = _.extend({}, data);
        transformed.ProductId = parseInt(data.ProductId);
        transformed.UserId = parseInt(data.UserId);
        transformed.Quantity = parseInt(data.Quantity);
        transformed.PricePerUnit = parseInt(data.PricePerUnit);
        transformed.TotalPrice = parseFloat(data.TotalPrice);
        console.log('transformed:', transformed);
        return JSON.stringify(transformed);
      },
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function() {
      $location.path('/coworking');
    })
    .error(function(data) {
      toastr.error('Error while trying to save changes');
    });
  };

  // Init scope variables
  var pickadateOptions = {
    format: 'yyyy-mm-dd'
  };
  $('.datepicker').pickadate(pickadateOptions);

  api.loadUsers(function(userData) {
    $scope.users = userData.users;
    $scope.usersById = userData.usersById;
    loadProducts();
  });

}]); // app.controller

})(); // closure

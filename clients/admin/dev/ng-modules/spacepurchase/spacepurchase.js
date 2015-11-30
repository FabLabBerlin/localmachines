(function(){

'use strict';
var app = angular.module('fabsmith.admin.space.purchase', ['ngRoute', 'ngCookies', 'fabsmith.admin.randomtoken', 'fabsmith.admin.api']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/space_purchases/:id', {
    templateUrl: 'ng-modules/spacepurchase/spacepurchase.html',
    controller: 'SpacePurchaseCtrl'
  });
}]); // app.config

app.controller('SpacePurchaseCtrl',
 ['$scope', '$routeParams', '$http', '$location', 'randomToken', 'api',
 function($scope, $routeParams, $http, $location, randomToken, api) {

  $scope.purchases = [];
  $scope.spacePurchase = {
    Id: $routeParams.id
  };
  $scope.spacesById = {};
  $scope.users = [];
  $scope.usersById = {};

  function parseInputTimes() {
    var sp = $scope.spacePurchase;
    sp.TimeStart = moment.tz(sp.DateStartLocal + ' ' + sp.TimeStartLocal, 'Europe/Berlin').toDate();
    sp.TimeEnd = moment.tz(sp.DateEndLocal + ' ' + sp.TimeEndLocal, 'Europe/Berlin').toDate();
  }

  function calculateQuantity() {
    console.log('$scope.timeChange()');
    parseInputTimes();
    var start = moment($scope.spacePurchase.TimeStart);
    var end = moment($scope.spacePurchase.TimeEnd);
    var duration = end.unix() - start.unix();
    console.log('duration=', duration);
    var quantity;
    switch ($scope.spacePurchase.PriceUnit) {
    case 'minute':
      quantity = duration / 60;
      break;
    case 'hour':
      quantity = duration / 3600;
      break;
    case 'day':
      quantity = duration / 24 / 3600;
      break;
    default:
      return;
    }
    $scope.spacePurchase.Quantity = quantity;
  }

  function calculateTotalPrice() {
    var totalPrice = $scope.spacePurchase.Quantity * $scope.spacePurchase.PricePerUnit;
    $scope.spacePurchase.TotalPrice = totalPrice.toFixed(2);
  }

  $scope.spaceChange = function() {
    console.log('$scope.spaceChange()');
    var spaceId = parseInt($scope.spacePurchase.ProductId);
    var space = $scope.spacesById[spaceId];
    $scope.spacePurchase.PricePerUnit = space.Product.Price;
    $scope.spacePurchase.PriceUnit = space.Product.PriceUnit;
    calculateQuantity();
    calculateTotalPrice();    
  };

  $scope.timeChange = function() {
    calculateQuantity();
    calculateTotalPrice();
  };

  $scope.priceUnitChange = function() {
    calculateQuantity();
    calculateTotalPrice();
  };

  $scope.updateSpacePurchase = function() {
    $http({
      method: 'PUT',
      url: '/api/purchases/' + $scope.spacePurchase.Id + '?type=space',
      headers: {'Content-Type': 'application/json' },
      data: $scope.spacePurchase,
      transformRequest: function(data) {
        return JSON.stringify(data);
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

  $scope.save = function() {
    parseInputTimes();
    $http({
      method: 'PUT',
      url: '/api/purchases/' + $scope.spacePurchase.Id + '?type=space',
      headers: {'Content-Type': 'application/json' },
      data: $scope.spacePurchase,
      transformRequest: function(data) {
        var transformed = _.extend({}, data);
        transformed.ProductId = parseInt(data.ProductId);
        transformed.UserId = parseInt(data.UserId);
        transformed.Quantity = parseFloat(data.Quantity);
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
      toastr.success('Space purchase updated');
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
    api.loadSpaces(function(spacesData) {
      $scope.spaces = spacesData.spaces;
      $scope.spacesById = spacesData.spacesById;
      api.loadSpacePurchase($scope.spacePurchase.Id, function(spacePurchase) {
        $scope.spacePurchase = spacePurchase;
        calculateTotalPrice();
        setTimeout(function() {
          $('.selectpicker').selectpicker('refresh');
        }, 100);
      });
    });
  });

}]); // app.controller

})(); // closure

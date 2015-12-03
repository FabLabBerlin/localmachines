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

  $scope.spaceChange = function() {
    $('#sp-space').selectpicker('refresh');
    console.log('$scope.spaceChange()');
    var spaceId = parseInt($scope.spacePurchase.ProductId);
    var space = $scope.spacesById[spaceId];
    $scope.spacePurchase.PricePerUnit = space.Product.Price;
    $scope.spacePurchase.PriceUnit = space.Product.PriceUnit;
    api.purchase.calculateQuantity($scope.spacePurchase);
    api.purchase.calculateTotalPrice($scope.spacePurchase);
  };

  $scope.userChange = function() {
    $('#sp-user').selectpicker('refresh');
  };

  $scope.timeChange = function() {
    api.purchase.calculateQuantity($scope.spacePurchase);
    api.purchase.calculateTotalPrice($scope.spacePurchase);
  };

  $scope.priceUnitChange = function() {
    api.purchase.calculateQuantity($scope.spacePurchase);
    api.purchase.calculateTotalPrice($scope.spacePurchase);
  };

  $scope.save = function() {
    api.purchase.parseInputTimes($scope.spacePurchase);
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
        api.purchase.calculateTotalPrice($scope.spacePurchase);
        setTimeout(function() {
          $('.selectpicker').selectpicker('refresh');
        }, 100);
      });
    });
  });

}]); // app.controller

})(); // closure

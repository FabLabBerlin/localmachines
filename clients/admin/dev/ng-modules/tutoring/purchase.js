(function(){

'use strict';

var app = angular.module('fabsmith.admin.tutoring.purchase',
  ['ngRoute', 'ngCookies', 'fabsmith.admin.api']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/tutoring/purchases/:id', {
    templateUrl: 'ng-modules/tutoring/purchase.html',
    controller: 'PurchaseCtrl'
  });
}]); // app.config

app.controller('PurchaseCtrl', ['$scope', '$routeParams', '$http', '$location', 'api',
  function($scope, $routeParams, $http, $location, api) {

  $scope.purchase = {
    Id: $routeParams.id
  };

  api.loadTutors(function(tutorData) {
    $scope.tutors = tutorData.tutors;
    $scope.tutorsById = tutorData.tutorsById;
    console.log('$scope.tutors:', $scope.tutors);
    api.loadUsers(function(userData) {
      $scope.users = userData.users;
      $scope.usersById = userData.usersById;
      api.loadTutoringPurchase($scope.purchase.Id, function(purchase) {
        $scope.purchase = purchase;
        console.log('$scope.purchase = ', $scope.purchase);
        setTimeout(function() {
          $('.selectpicker').selectpicker('refresh');
        }, 100);
      });
    });
  });

  // Init scope variables
  var pickadateOptions = {
    format: 'yyyy-mm-dd'
  };
  $('.datepicker').pickadate(pickadateOptions);

  $scope.tutorChange = function() {
    var tutorId = parseInt($scope.purchase.ProductId);
    var tutor = $scope.tutorsById[tutorId];
    $scope.purchase.PricePerUnit = tutor.Price;
    $scope.purchase.PriceUnit = tutor.PriceUnit;
    api.purchase.calculateQuantity($scope.purchase);
    api.purchase.calculateTotalPrice($scope.purchase);
  };

  $scope.timeChange = function() {
    api.purchase.calculateQuantity($scope.purchase);
    api.purchase.calculateTotalPrice($scope.purchase);
  };

  $scope.priceUnitChange = function() {
    api.purchase.calculateQuantity($scope.purchase);
    api.purchase.calculateTotalPrice($scope.purchase);
  };


  $scope.save = function() {
    api.purchase.parseInputTimes($scope.purchase);
    $http({
      method: 'PUT',
      url: '/api/purchases/' + $scope.purchase.Id + '?type=tutor',
      headers: {'Content-Type': 'application/json' },
      data: $scope.purchase,
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
      toastr.success('Tutoring purchase updated');
    })
    .error(function(data) {
      toastr.error('Error while trying to save changes');
    });
  };

  $scope.cancel = function() {
    $location.path('/tutoring');
  };

}]); // app.controller

})(); // closure
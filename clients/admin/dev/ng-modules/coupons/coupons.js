(function(){

'use strict';

var app = angular.module('fabsmith.admin.coupons', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/coupons', {
    templateUrl: 'ng-modules/coupons/coupons.html',
    controller: 'CouponsCtrl'
  });
}]); // app.config

app.controller('CouponsCtrl',
 ['$scope', '$http', '$location', '$cookies', '$q', 'api',
 function($scope, $http, $location, $cookies, $q, api) {

  function loadCoupons() {
    $http({
      url: '/api/coupons',
      params: {
        location: $cookies.locationId
      }
    })
    .success(function(coupons) {
      $scope.coupons = coupons;
    })
    .error(function() {
      toastr.error('Failed to get coupons');
    });
  }

  $scope.generate = function() {
    api.prompt('Do you really want to create ' + $scope.n + ' coupons? ', function() {
      $http({
        method: 'POST',
        url: '/api/coupons',
        params: {
          location: $cookies.locationId
        },
        data: {
          static_code: $scope.staticCode,
          n: $scope.n,
          value: $scope.value
        }
      })
      .success(function() {
        loadCoupons();
      })
      .error(function() {
        toastr.error('Failed to generate coupons.  Please try again later.');
      });
    });
  };

  loadCoupons();

}]);

})(); // closure

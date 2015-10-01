(function(){

'use strict';

var app = angular.module('fabsmith.admin.reservations', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/reservations', {
    templateUrl: 'ng-modules/reservations/reservations.html',
    controller: 'ReservationsCtrl'
  });
}]); // app.config

app.controller('ReservationsCtrl', ['$scope', '$http', '$location', '$cookieStore', 
 function($scope, $http, $location, $cookieStore) {

  $scope.reservations = [];

  $http({
    method: 'GET',
    url: '/api/reservations',
    params: {
      ac: new Date().getTime()
    }
  })
  .success(function(data) {
    $scope.reservations = data;
  })
  .error(function() {
    toastr.error('Failed to get reservations');
  });

}]); // app.controller

})(); // closure
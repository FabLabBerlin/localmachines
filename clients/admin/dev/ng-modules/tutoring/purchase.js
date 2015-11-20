(function(){

'use strict';

var app = angular.module('fabsmith.admin.tutoring.purchase', 
  ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/tutoring/purchase', {
    templateUrl: 'ng-modules/tutoring/purchase.html',
    controller: 'PurchaseCtrl'
  });
}]); // app.config

app.controller('PurchaseCtrl', ['$scope', '$http', '$location', 
  function($scope, $http, $location) {

  $scope.user = {
    Name: 'Sugru Meyer'
  };

  $scope.purchase = {
    StartTimeDate: '23 Nov 15',
    StartTimeTime: '15:00',
    EndTimeDate: '23 Nov 15',
    EndTimeTime: '17:00',
    TimeReserved: '2h 0m',
    TimeTimed: '1h 12m',
    PriceTotal: '120.00'
  };

  $scope.save = function() {
    toastr.success('Tutoring purchase saved');
    $location.path('/tutoring');
  };

  $scope.cancel = function() {
    $location.path('/tutoring');
  };

}]); // app.controller

})(); // closure
(function(){

'use strict';

var app = angular.module('fabsmith.admin.tutoring.purchase', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/tutoring/purchase', {
    templateUrl: 'ng-modules/tutoring/purchase.html',
    controller: 'PurchaseCtrl'
  });
}]); // app.config

app.controller('PurchaseCtrl', ['$scope', '$http', '$location', 
  function($scope, $http, $location) {

  $scope.purchase = {

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
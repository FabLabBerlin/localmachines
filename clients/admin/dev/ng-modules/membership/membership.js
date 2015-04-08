(function(){

'use strict';

var app = angular.module('fabsmith.admin.membership', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/membership/:membershipId', {
    templateUrl: 'ng-modules/membership/membership.html',
    controller: 'MembershipCtrl'
  });
}]); // app.config

app.controller('MembershipCtrl', ['$scope', '$http', '$location', '$filter', 
 function($scope, $http, $location, $filter) {

  $scope.membership = {
    Id: 1,
    Title: 'My Membership',
    ShortName: 'MM',
    Duration: 1,
    Unit: 'day',
    Price: 30,
    MachinePriceDeduction: 50
  };

  $scope.membership.Price = $filter('currency')($scope.membership.Price,'',2);
  
  

}]); // app.controller

})(); // closure
(function(){

'use strict';

var app = angular.module('fabsmith.admin.membership', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/membership/:membershipId', {
    templateUrl: 'ng-modules/membership/membership.html',
    controller: 'MembershipCtrl'
  });
}]); // app.config

app.controller('MembershipCtrl', 
 ['$scope', '$http', '$location', '$filter', '$routeParams',
 function($scope, $http, $location, $filter, $routeParams) {

  $scope.membership = {
    Id: $routeParams.membershipId
  };

  // Load membership
  $http({
    method: 'GET',
      url: '/api/memberships/' + $scope.membership.Id,
      params: {
        anticache: new Date().getTime()
      }
  })
  .success(function(membershipModel) {
    $scope.membership = membershipModel;
    $scope.membership.Price = $filter('currency')($scope.membership.Price,'',2);
  })
  .error(function() {
    toastr.error('Failed to get membership');
  });

  /*
  $scope.membership = {
    Id: 1,
    Title: 'My Membership',
    ShortName: 'MM',
    Duration: 1,
    Unit: 'day',
    Price: 30,
    MachinePriceDeduction: 50
  };
  */

}]); // app.controller

})(); // closure
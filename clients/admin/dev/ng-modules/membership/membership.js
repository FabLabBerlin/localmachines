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

  $scope.machines = [];
  $scope.membership = {
    Id: $routeParams.membershipId
  };

  // Load machines first
  $http({
    method: 'GET',
    url: '/api/machines',
    params: {
      anticache: new Date().getTime()
    }
  })
  .success(function(machines) {
    $scope.machines = machines;
    $scope.loadMembership();
  })
  .error(function() {
    toastr.error('Failed to load machines');
  });

  // Load membership
  $scope.loadMembership = function() {
    $http({
      method: 'GET',
      url: '/api/memberships/' + $scope.membership.Id,
      params: {
        anticache: new Date().getTime()
      }
    })
    .success(function(membershipModel) {
      $scope.membership = membershipModel;
      
      // Parse affected machines JSON as it is passed here as string
      $scope.membership.AffectedMachines = 
       $.parseJSON($scope.membership.AffectedMachines);

      // Search for machines with the same IDs as the AffectedMachines
      // and set them as checked
      for (var i = 0; i < $scope.membership.AffectedMachines.length; i++) {
        for (var j = 0; j < $scope.machines.length; j++) {
          if (parseInt($scope.machines[j].Id) === 
              parseInt($scope.membership.AffectedMachines[i])) {
            $scope.machines[j].Checked = true;
          }
        }
      }

      // Filter as currency
      $scope.membership.Price = 
       $filter('currency')($scope.membership.Price,'',2);
    })
    .error(function() {
      toastr.error('Failed to get membership');
    });
  };

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

  $scope.updateMembership = function() {

  };

}]); // app.controller

})(); // closure
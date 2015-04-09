(function(){

'use strict';

var app = angular.module('fabsmith.admin.memberships', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/memberships', {
    templateUrl: 'ng-modules/memberships/memberships.html',
    controller: 'MembershipsCtrl'
  });
}]); // app.config

app.controller('MembershipsCtrl', ['$scope', '$http', '$location', 
 function($scope, $http, $location) {

  // Load all memberships
  $http({
    method: 'GET',
    url: '/api/memberships',
    params: {
      anticache: new Date().getTime()
    }
  })
  .success(function(data) {
    $scope.memberships = data;
  })
  .error(function() {
    toastr.error('Failed to load memberships');
  });

  /*
  $scope.memberships = [
    {
      Id: 1,
      Title: 'My Membership',
      ShortName: 'MM',
      Duration: 1,
      Unit: 'day',
      Price: 30,
      MachinePriceDeduction: 50
    }
  ];
  */

  $scope.addMembershipPrompt = function() {
    vex.dialog.prompt({
      message: 'Enter membership name',
      placeholder: 'Membership name',
      callback: $scope.membershipPromptCallback
    });
  };

  $scope.membershipPromptCallback = function(value) {
    if (value) {    
      $scope.addMembership(value);
    } else if (value !== false) {
      toastr.error('No membership name');
    }
  };

  $scope.addMembership = function(membershipName) {
    $http({
      method: 'POST',
      url: '/api/memberships',
      params: {
        mname: membershipName,
        anticache: new Date().getTime()
      }
    })
    .success(function(membershipId) {
      $scope.editMembership(membershipId);
    })
    .error(function() {
      toastr.error('Failed to create membership');
    });
  };

  $scope.editMembership = function(membershipId) {
    $location.path('/membership/' + membershipId);
  };

}]); // app.controller

})(); // closure
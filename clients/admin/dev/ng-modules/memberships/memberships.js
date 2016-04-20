(function(){

'use strict';

var app = angular.module('fabsmith.admin.memberships', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/memberships', {
    templateUrl: 'ng-modules/memberships/memberships.html',
    controller: 'MembershipsCtrl'
  });
}]); // app.config

app.controller('MembershipsCtrl',
 ['$scope', '$http', '$location', '$cookies',
 function($scope, $http, $location, $cookies) {

  // Load all memberships
  $http({
    method: 'GET',
    url: '/api/memberships',
    params: {
      location: $cookies.locationId,
      ac: new Date().getTime()
    }
  })
  .success(function(data) {
    $scope.memberships = data;
  })
  .error(function() {
    toastr.error('Failed to load memberships');
  });

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
        location: $cookies.locationId,
        mname: membershipName,
        ac: new Date().getTime()
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

  $scope.setShowArchived = function(show) {
    $scope.showArchived = show;
  };

}]); // app.controller

app.filter('membershipsFilter', function() {
  return function(memberships, scope) {
    return _.filter(memberships, function(membership) {
      return scope.showArchived || !membership.Archived;
    });
  }
});

})(); // closure
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

  $scope.addMembershipPrompt = function() {
    vex.dialog.prompt({
      message: 'Enter membership name',
      placeholder: 'Membership name',
      callback: $scope.membershipPromptCallback
    });
  };

  $scope.membershipPromptCallback = function(value) {
    console.log(value);
  };

  $scope.addMembership = function() {
    toastr.warning('adding membership');
  };

  $scope.editMembership = function() {
    toastr.warning('editing membership');
  };

}]); // app.controller

})(); // closure
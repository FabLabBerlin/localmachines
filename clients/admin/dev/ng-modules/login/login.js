(function(){

'use strict';

var app = angular.module('fabsmith.admin.login', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/login', {
    templateUrl: 'ng-modules/login/login.html',
    controller: 'LoginCtrl'
  });
}]); // app.config

app.controller('LoginCtrl',
 ['$rootScope', '$scope', '$http', '$location', '$cookies',
 function($rootScope, $scope, $http, $location, $cookies) {
  $scope.getUserData = function(userId) {
    $http({
      method: 'GET',
      url: '/api/users/' + userId,
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data){
      $rootScope.mainMenu.userFullName = data.FirstName + ' ' + data.LastName;
      $location.path('/machines');
    })
    .error(function(data, status){
      console.log('Status: ' + status);
      console.log('Data' + data);
      toastr.error('Could not get user data');
    });
  };

  $scope.login = function() {
    var locationId = $('select[name="location"]').val();
    locationId = parseInt(locationId);

    $http({
      method: 'POST',
      url: '/api/users/login',
      data: {
        username: $scope.username,
        password: $scope.password,
        location: locationId,
        admin: true
      },
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      if (data.UserId) {
        $scope.getUserData(data.UserId);
      }
    })
    .error(function() {
      toastr.error('Failed to log in');
    });
  };

  $scope.getLocations = function() {
    $http({
      method: 'GET',
      url: '/api/locations',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(locations) {
      $scope.locations = locations;
    })
    .error(function() {
      toastr.error('Failed to load locations');
    });
  };

  $scope.getLocations();
  
}]); // app.controller

})(); // closure
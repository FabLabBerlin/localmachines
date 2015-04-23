(function(){

'use strict';

angular.module('fabsmith.login', ['ngRoute', 'ngCookies'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/login', {
    templateUrl: 'ng-modules/login/login.html',
    controller: 'LoginCtrl'
  });
}])

.controller('LoginCtrl', ['$scope', '$http', '$location', function($scope, $http, $location) {
  
  // Local login function - if we do it by entering username and 
  // password in the browser
  $scope.login = function() {
    // Attempt to login via API
    $http({
      method: 'POST',
      url: '/api/users/login',
      data: {
        username: $scope.username,
        password: $scope.password
      },
      params: { ac: new Date().now() }
    })
    .success(function(loginResponse) {
      if (loginResponse.UserId) {
        $scope.getUser(loginResponse.UserId);
      } else {
        toastr.error('Failed to log in');
      }
    })
    .error(function(data, status) {
      toastr.error('Failed to log in');
    });
  };

  $scope.getUser = function(userId) {
    $http({
      method: 'GET',
      url: '/api/users/' + userId,
      params: { ac: new Date().now() }
    })
    .success(function(userData){
      $scope.onUserDataLoaded(userData);
    })
    .error(function(data, status){
      toastr.error('Could not get user data');
    });
  };

  $scope.onUserDataLoaded = function(userData){
    $scope.$emit('user-login', userData);
    $location.path('/machines');
  };

  // Make the main controller scope accessible from outside
  // So we can use our Android app to call login function
  window.LOGIN_CTRL_SCOPE = $scope;

  // Call this from Android app as LOGIN_CTRL_SCOPE.login("user", "pass");
  $scope.androidLogin = function(username, password) {
    $scope.username = username;
    $scope.password = password;
    $scope.login();
  };

}]);

})();
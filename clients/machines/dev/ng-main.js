(function(){

'use strict';

// Declare app level module which depends on views, and components
var app = angular.module('fabsmith', [
  'ngRoute',
  'ngCookies',
  'fabsmith.login',
  'fabsmith.machines',
  'fabsmith.logout',
  'fabsmith.version'
]);

// This checks whether an user is logged in always before 
// switching to a new view
app.run(['$rootScope', '$location', '$http', 
  function($rootScope, $location, $http) {

  // On each location change
  $rootScope.$on('$locationChangeStart', function(event, newUrl, oldUrl) {

    // Get requested angular path
    var newPath = newUrl.split('#')[1];
    
    // If it is not login (main) view, 
    // check if the user is logged in
    if (newPath !== '/login') {
      $http({
        method: 'POST',
        url: '/api/users/login',
        params: {
          username: 'blank', // TODO: randomize?
          password: 'blank'
        }
      })
      .success(function(data) {
        if (data.Status !== 'logged') {
          event.preventDefault();
          $location.path('/login');
        } else {
          if (newPath) {
            $location.path(newPath);
          } else {
            $location.path('/machines');
          }
        }
      })
      .error(function() {
        event.preventDefault();
        $location.path('/login');
      });
    }

  });
}]);

// Configure http provider to send regular form POST data instead of JSON
app.config(['$httpProvider', function($httpProvider) {

  $httpProvider.defaults.transformRequest = function(data){
    if (data === undefined) {
      return data;
    }
    return $.param(data);
  };
  $httpProvider.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded; charset=UTF-8';
}]);

// Main controller, checks if user logged in
app.controller('MainCtrl', ['$scope', '$http', '$location', '$cookieStore', '$cookies', 
function($scope, $http, $location, $cookieStore, $cookies){

  // Configure toastr default location
  toastr.options.positionClass = 'toast-bottom-left';

  // Configure vex theme
  vex.defaultOptions.className = 'vex-theme-custom';

  // Configure root scope so Android can access
  window.ROOT_SCOPE = $scope;

  // Store user data on user login
  $scope.putUserData = function(data) {
    for (var key in data) {
      if (data.hasOwnProperty(key)) {
        $cookieStore.put(key, data[key]);
      }
    }
  };
  $scope.$on('user-login', function (event, data){
    $scope.putUserData(data);
  });

  // Clear user data on user logout
  $scope.removeUserData = function() {
    for (var key in $cookies) {
      if ($cookies.hasOwnProperty(key)) {
        $cookieStore.remove(key);
      }
    }
  };

  // Log out on logout event
  $scope.logout = function() {
    $scope.removeUserData();

    $http({
      method: 'GET',
      url: '/api/users/logout'
    })
    .success(function() {
      $location.path('/');
    })
    .error(function() {
      toastr.error('Failed to log out. Probably server down.');
      $location.path('/');
    });

  };
  $scope.$on('logout', function (event, data){
    $scope.logout();
  });

}]);

})(); // closure

(function(){

'use strict';

// Declare app level module which depends on views, and components
var app = angular.module('fabsmith', [
  'ngRoute',
  'fabsmith.admin.login',
  'fabsmith.admin.dashboard',
  'fabsmith.admin.mainmenu',
  'fabsmith.admin.user',
  'fabsmith.admin.activations',
  'fabsmith.admin.machines',
  'fabsmith.admin.machine',
  'fabsmith.admin.memberships',
  'fabsmith.admin.membership',
  'fabsmith.admin.bookings',
  'fabsmith.admin.invoices',
  'fabsmith.admin.users',
  'fabsmith.admin.randomtoken',
  'fabsmith.version'
]);

// This checks whether an user is logged in always before switching to a new view
app.run(['$rootScope', '$location', '$http', '$cookieStore', 
 function($rootScope, $location, $http, $cookieStore) {
  
  $rootScope.$on('$locationChangeStart', 
   function(event, newUrl, oldUrl) {

    // Get requested angular path
    var newPath = newUrl.split('#')[1];
    var userId = $cookieStore.get('Id');
    
    // Only check if user logged in if requested page is other
    // than the login page
    if (newPath !== '/login') {

      // Just use the cookie user ID to check log in status.
      // Very simple and no double requests.
      if (userId) {
        if (newPath) {
          $location.path(newPath);
        } else {
          $location.path('/dashboard');
        }
      } else {
        $location.path('/login');
      }

    }

  });
}]); // app.run

// Configure http provider to send regular form POST data instead of JSON
app.config(['$httpProvider', function($httpProvider) {
  $httpProvider.defaults.transformRequest = function(data){
    if (data === undefined) {
        return data;
    }
    return $.param(data);
  };
  $httpProvider.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded; charset=UTF-8';

  // Disable caching (http://stackoverflow.com/a/19771501/485185)
  if (!$httpProvider.defaults.headers.get) {
    $httpProvider.defaults.headers.get = {};
  }
  //disable IE ajax request caching
  $httpProvider.defaults.headers.get['If-Modified-Since'] = 'Mon, 26 Jul 1997 05:00:00 GMT';
  // extra
  $httpProvider.defaults.headers.get['Cache-Control'] = 'no-cache';
  $httpProvider.defaults.headers.get.Pragma = 'no-cache';

}]); // app.config

// Main controller, checks if user logged in
app.controller('MainCtrl', ['$scope', '$http', '$location', '$cookieStore', '$cookies', 
 function($scope, $http, $location, $cookieStore, $cookies){

  $scope.mainMenu = {visible: false};

  // Configure toastr default location
  toastr.options.positionClass = 'toast-bottom-left';

  // Configure vex theme
  vex.defaultOptions.className = 'vex-theme-plain';

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
    $scope.mainMenu.visible = true;
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
    $scope.mainMenu.visible = false;

    $http({
      method: 'GET',
      url: '/api/users/logout',
      params: {
        anticache: new Date().getTime()
      }
    })
    .success(function() {
      $location.path('/');
    })
    .error(function() {
      alert('Failed to log out. Probably server down.');
      $location.path('/');
    });

  };
  $scope.$on('logout', function (event, data){
    $scope.logout();
  });
}]);

})(); // closure

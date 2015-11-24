(function(){

'use strict';

// Declare app level module which depends on views, and components
var app = angular.module('fabsmith', [
  'ngRoute',
  'fabsmith.admin.login',
  'fabsmith.admin.coworking',
  'fabsmith.admin.coworking.table',
  'fabsmith.admin.dashboard',
  'fabsmith.admin.mainmenu',
  'fabsmith.admin.user',
  'fabsmith.admin.activations',
  'fabsmith.admin.machines',
  'fabsmith.admin.machine',
  'fabsmith.admin.memberships',
  'fabsmith.admin.membership',
  'fabsmith.admin.productlist',
  'fabsmith.admin.reservation',
  'fabsmith.admin.reservations',
  'fabsmith.admin.invoices',
  'fabsmith.admin.users',
  'fabsmith.admin.randomtoken',
  'fabsmith.admin.settings',
  'fabsmith.version',
  'fabsmith.admin.tutoring',
  'fabsmith.admin.tutoring.tutor',
  'fabsmith.admin.tutoring.purchase',
  'fabsmith.admin.space',
  'fabsmith.admin.space.purchase',
  'fabsmith.admin.spaces',
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

      $http({
        method: 'POST',
        url: '/api/users/login',
        params: {
          username: 'blank',
          password: 'blank',
          ac: new Date().getTime()
        }
      })
      .success(function(data){
        if (data.Status !== 'logged') {
          $rootScope.mainMenu.visible = false;
          $location.path('/login');
        } else {
          $rootScope.mainMenu.visible = true;
          if (newPath) {
            $location.path(newPath);
          } else {
            $location.path('/dashboard');
          }
        }
      })
      .error(function(){
        $rootScope.mainMenu.visible = false;
        $location.path('/login');
      });

    } // if newPath

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
}]); // app.config

// Main controller, checks if user logged in
app.controller('MainCtrl', 
 ['$scope', '$http', '$location', '$cookieStore', '$cookies', '$rootScope', 
 function($scope, $http, $location, $cookieStore, $cookies, $rootScope){

  $rootScope.mainMenu = {visible: false};

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
      params: { ac: new Date().getTime() }
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

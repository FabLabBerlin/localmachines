(function(){

'use strict';

// Declare app level module which depends on views, and components
var app = angular.module('fabsmith', [
  'ngRoute',
  'fabsmith.admin.login',
  'fabsmith.admin.activation',
  'fabsmith.admin.activations',
  'fabsmith.admin.api',
  'fabsmith.admin.coworking',
  'fabsmith.admin.coworking.purchase',
  'fabsmith.admin.coworking.product',
  'fabsmith.admin.dashboard',
  'fabsmith.admin.mainmenu',
  'fabsmith.admin.user',
  'fabsmith.admin.machines',
  'fabsmith.admin.machine',
  'fabsmith.admin.memberships',
  'fabsmith.admin.membership',
  'fabsmith.admin.priceunit',
  'fabsmith.admin.productlist',
  'fabsmith.admin.reservation',
  'fabsmith.admin.reservations',
  'fabsmith.admin.reservations.toggle',
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
app.run(['$rootScope', '$location', '$http', '$cookies', '$q',
 function($rootScope, $location, $http, $cookies, $q) {
  
  $rootScope.$on('$locationChangeStart', 
   function(event, newUrl, oldUrl) {

    // Get requested angular path
    var newPath = newUrl.split('#')[1];
    
    // Only check if user logged in if requested page is other
    // than the login page
    if (newPath !== '/login') {

      var getLocationPromise = $http({
        method: 'GET',
        url: '/api/locations/' + $cookies.locationId,
        params: {
          location: $cookies.locationId,
          ac: new Date().getTime()
        }
      });

      var getUserPromise = $http({
        method: 'GET',
        url: '/api/users/current',
        params: {
          ac: new Date().getTime()
        }
      });

      $q.all([
        getLocationPromise,
        getUserPromise
      ])
      .then(function(results) {
        var location = results[0].data;
        var user = results[1].data;

        $rootScope.mainMenu.visible = true;
        $rootScope.mainMenu.userFullName = user.FirstName + ' ' + user.LastName;
        $rootScope.mainMenu.location = location;

        if (newPath) {
          $location.path(newPath);
        } else {
          $location.path('/machines');
        }
      })
      .catch(function(err) {
        console.log('err:', err);
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
 ['$scope', '$http', '$location', '$cookies', '$rootScope', 
 function($scope, $http, $location, $cookies, $rootScope){

  $rootScope.mainMenu = {visible: false};

  // Configure toastr default location
  toastr.options.positionClass = 'toast-bottom-left';

  // Configure vex theme
  vex.defaultOptions.className = 'vex-theme-plain';

  // Log out on logout event
  $scope.logout = function() {
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

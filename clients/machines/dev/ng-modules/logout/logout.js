(function(){

'use strict';

angular.module('fabsmith.logout', ['ngRoute'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/logout', {
    templateUrl: 'ng-modules/logout/logout.html',
    controller: 'LogoutCtrl'
  });
}])

.controller('LogoutCtrl', ['$scope', '$http', '$location', function($scope, $http, $location) {

  // Activate countdown
  $scope.abortLogout = function() {
    $location.path('/machines');
  };

  $scope.logout = function() {
    $http({
      method: 'GET',
      url: '/api/users/logout'
    })
    .success(function() {
      $location.path('/');
    })
    .error(function() {
      toastr.error('Failed to log out. Probably server down.');
    });
  };

}]);

})(); // closure
(function(){

'use strict';

var app = angular.module('fabsmith.admin.users', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/users', {
    templateUrl: 'ng-modules/users/users.html',
    controller: 'UsersCtrl'
  });
}]); // app.config

app.controller('UsersCtrl',
 ['$scope', '$http', '$location', '$cookies', '$q',
 function($scope, $http, $location, $cookies, $q) {

  $scope.users = [];

  $scope.getAllUsers = function() {
    $q.all([
      $http({
        method: 'GET',
        url: '/api/users',
        params: {
          ac: new Date().getTime(),
          location: $cookies.get('location')
        }
      }),
      $http({
        method: 'GET',
        url: '/api/user_locations',
        params: {
          ac: new Date().getTime(),
          location: $cookies.get('location')
        }
      })
    ])
    .then(function(result) {
      var users = result[0].data;
      var userLocations = result[1].data;

      $scope.users = users;
      $scope.usersById = {};
      _.each(users, function(user) {
        $scope.usersById[user.Id] = user;
      });
      _.each(userLocations, function(userLocation) {
        var u = $scope.usersById[userLocation.UserId];
        if (u) {
          u.UserLocation = userLocation;
        }
      });
    })
    .catch(function(data, status) {
      toastr.error('Failed to get all users');
    });
  };

  $scope.getAllUsers();

  $scope.addUserPrompt = function() {
    vex.dialog.prompt({
      message: 'Enter user email',
      placeholder: 'user@example.com',
      callback: $scope.addUserPromptCallback
    });
  };

  $scope.addUserPromptCallback = function(email) {
    if (email) {
      // TODO: validate email
      $scope.addUser(email);
    }
  };

  $scope.addUser = function(email) {
    $http({
      method: 'POST',
      url: '/api/users',
      data: {
        email: email,
        location: $cookies.get('location')
      },
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(user) {
      toastr.info('New user created');
      $location.path('/users/' + user.Id);
    })
    .error(function() {
      toastr.error('Error while trying to create new user');
    });
  };

  $scope.editUser = function(userId) {
    $location.path('/users/' + userId);
  };

  $scope.setShowArchived = function(show) {
    $scope.showArchived = show;
  };

}]); // app.controller

app.filter('usersFilter', function() {
  return function(users, scope) {
    return _.filter(users, function(user) {
      var userArchived = user.UserLocation &&
        user.UserLocation.UserRole === 'archived';
      return scope.showArchived || !userArchived;
    });
  };
});

app.directive('userListItem', ['$location', function($location) {
  return {
    templateUrl: 'ng-modules/users/user-list-item.html',
    restrict: 'E',
    controller: ['$scope', '$element', function($scope, $element) {

    }]
  };
}]);

app.directive('userListHead', ['$location', function($location) {
  return {
    templateUrl: 'ng-modules/users/user-list-head.html',
    restrict: 'E'
  };
}]);

})(); // closure
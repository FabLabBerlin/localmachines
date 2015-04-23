(function(){

'use strict';

var app = angular.module('fabsmith.admin.users', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/users', {
    templateUrl: 'ng-modules/users/users.html',
    controller: 'UsersCtrl'
  });
}]); // app.config

app.controller('UsersCtrl', ['$scope', '$http', '$location', 
 function($scope, $http, $location) {

  $scope.users = [];

  $scope.getAllUsers = function() {
    $http({
      method: 'GET',
      url: '/api/users',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(users) {
      $scope.users = users;
    })
    .error(function(data, status) {
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
        email: email
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
}]); // app.controller

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
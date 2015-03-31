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

    // Get all users
    $http({
        method: 'GET',
        url: '/api/users',
        params: {
            anticache: new Date().getTime()
        }
    })
    .success(function(data) {
        $scope.users = data;
    })
    .error(function(data, status) {
        toastr.error('Failed to get all users');
    });

    $scope.addUser = function() {
        var email = prompt('Please enter E-Mail for new user:');
        if (email) {
            $http({
                method: 'POST',
                url: '/api/users',
                data: {
                    email: email,
                    anticache: new Date().getTime()
                }
            })
            .success(function(data) {
                toastr.info('New user created');
                $location.path('/users/' + data.Id);
            })
            .error(function() {
                toastr.error('Error while trying to create new user');
            });
        }
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
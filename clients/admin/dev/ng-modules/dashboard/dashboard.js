(function(){

'use strict';

var app = angular.module('fabsmith.backoffice.dashboard', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/dashboard', {
    templateUrl: 'ng-modules/dashboard/dashboard.html',
    controller: 'DashboardCtrl'
  });
}]); // app.config

app.controller('DashboardCtrl', ['$scope', '$http', function($scope, $http) {
	$scope.testVar = 'This is a test variable';

	$scope.users = [];

	// Get current user machines
	$http({
		method: 'GET',
		url: '/api/users',
		params: {
			anticache: new Date().getTime()
		}
	})
	.success(function(data) {
		if (data.Status === 'error') {
			alert('msg: ' + data.Users);
		} else if (data.Users.length > 0) {
			$scope.users = data.Users;
		} else {
			alert('Error loading users');
		}
	})
	.error(function() {
		alert('Error loading users');
	});

	$scope.addUser = function() {
		alert('Add User!');
	};
}]); // app.controller

app.directive('fsUserItem', ['$location', function($location) {
	return {
		templateUrl: 'ng-modules/dashboard/user-item.html',
		restrict: 'E',
		controller: ['$scope', '$element', function($scope, $element) {
			$scope.editUser = function(userId) {
				$location.path('/users/' + userId);
			};
		}]
	};
}]);

})(); // closure
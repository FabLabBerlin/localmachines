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
		console.log('Got all users');
		console.log(data);
		$scope.users = data;
	})
	.error(function(data, status) {
		console.log('Could not get users');
		console.log('Data: ' + data);
		console.log('Status code: ' + status);
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
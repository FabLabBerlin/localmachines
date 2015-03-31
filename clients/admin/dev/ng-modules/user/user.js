(function(){

'use strict';
var app = angular.module('fabsmith.admin.user', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
	$routeProvider.when('/users/:userId', {
		templateUrl: 'ng-modules/user/user.html',
		controller: 'UserCtrl'
	});
}]); // app.config

app.controller('UserCtrl', ['$scope', '$routeParams', '$http', '$location', function($scope, $routeParams, $http, $location) {
	$scope.user = {
		Id: $routeParams.userId
	};
	$scope.userMachines = [];
	$scope.userRoles = {};

	$http({
		method: 'GET',
		url: '/api/users/' + $scope.user.Id,
		params: {
			anticache: new Date().getTime()
		}
	})
	.success(function(data) {
		console.log('Got user');
		console.log(data);
		$scope.user = data;
	})
	.error(function(data, status) {
		console.log('Could not get user');
		console.log('Data: ' + data);
		console.log('Status code: ' + status);
	});

	$http({
		method: 'GET',
		url: '/api/users/' + $scope.user.Id + '/machines',
		params: {
			anticache: new Date().getTime()
		}
	})
	.success(function(data) {
		console.log('Got user machines');
		console.log(data);
		$scope.userMachines = data;
	})
	.error(function(data, status) {
		console.log('Could not get user machines');
		console.log('Data: ' + data);
		console.log('Status code: ' + status);
	});

	$http({
		method: 'GET',
		url: '/api/users/' + $scope.user.Id + '/roles',
		params: {
			anticache: new Date().getTime()
		}
	})
	.success(function(data) {
		console.log('Got user roles');
		console.log(data);
		$scope.userRoles = data;
	})
	.error(function(data, status) {
		console.log('Could not get user roles');
		console.log('Data: ' + data);
		console.log('Status code: ' + status);
	});

	$scope.cancel = function() {
		if (confirm('All changes will be discarded, click ok to continue.')) {
			$location.path('/users');
		}
	};

	$scope.deleteUser = function() {
		var email = prompt("Do you really want to delete this user? Please enter user's E-Mail address to continue");
		if (email === $scope.user.Email) {
			$http({
				method: 'DELETE',
				url: '/api/users/' + $scope.user.Id
			})
			.success(function(data) {
				toastr.info('User deleted');
				$location.path('/users');
			})
			.error(function() {
				toastr.error('Error while trying to delete user');
			});
		} else {
			toastr.warning('Delete User canceled.');
		}
	};

	$scope.saveUser = function() {
		console.log('user model:', $scope.user);
		$http({
			method: 'PUT',
			url: '/api/users/' + $scope.user.Id,
			headers: {'Content-Type': 'application/json' },
			data: {
				User: $scope.user,
				UserRoles: $scope.userRoles
			},
			transformRequest: function(data) {
				console.log('data to send:', data);
				return JSON.stringify(data);
			}
		})
		.success(function() {
			toastr.info('Changes saved.');
		})
		.error(function() {
			toastr.error('Error while trying to save changes');
		});
	};
}]); // app.controller

})(); // closure
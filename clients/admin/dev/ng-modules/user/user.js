(function(){

'use strict';
var app = angular.module('fabsmith.backoffice.user', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
	$routeProvider.when('/users/:userId', {
		templateUrl: 'ng-modules/user/user.html',
		controller: 'UserCtrl'
	});
}]); // app.config

app.controller('UserCtrl', ['$scope', '$routeParams', '$http', function($scope, $routeParams, $http) {
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
}]); // app.controller

})(); // closure
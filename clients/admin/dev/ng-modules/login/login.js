(function(){

'use strict';

var app = angular.module('fabsmith.admin.login', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/login', {
    templateUrl: 'ng-modules/login/login.html',
    controller: 'LoginCtrl'
  });
}]); // app.config

app.controller('LoginCtrl', ['$scope', '$http', '$location', function($scope, $http, $location) {
	// Local login function - if we do it by entering username and password in the browser
	$scope.login = function() {
		$http({
			method: 'POST',
			url: '/api/users/login',
			data: {
				username: $scope.username,
				password: $scope.password,
				anticache: new Date().getTime()
			}
		})
		.success(function(data) {

			if (data.UserId) {
				console.log('User ID: ' + data.UserId);

				// Get user data
				$http({
					method: 'GET',
					url: '/api/users/' + data.UserId
				})
				.success(function(data){
					console.log('Got user data');
					$scope.$emit('user-login', data);
					$location.path('/dashboard');
				})
				.error(function(data, status){
					console.log('Status: ' + status);
					console.log('Data' + data);
					toastr.error('Could not get user data');
				});
				
			} // if data.UserId
		})
		.error(function() {
			toastr.error('Failed to log in');
		});
	};
	
}]); // app.controller

})(); // closure
(function(){

'use strict';

var app = angular.module('fabsmith.backoffice.login', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/login', {
    templateUrl: '/static/dev/admin/login/login.html',
    controller: 'LoginCtrl'
  });
}]); // app.config

app.controller('LoginCtrl', ['$scope', '$http', '$location', function($scope, $http, $location) {
	// Local login function - if we do it by entering username and password in the browser
	$scope.login = function() {
		$http({
			method: 'POST',
			url: '/api/login',
			params: {
				username: $scope.username,
				password: md5($scope.password),
				anticache: new Date().getTime()
			}
		})
		.success(function(data) {
			if (data.Status === 'error') {
				alert(data.Message);
			} else if (data.Status === 'logged' || data.Status === 'ok'){
				$scope.$emit('user-login', data);
				$location.path('/dashboard');
			}
		})
		.error(function() {
			alert('Failed to log in');
		});
	}
}]); // app.controller

})(); // closure
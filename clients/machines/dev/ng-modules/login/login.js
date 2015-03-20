(function(){

'use strict';

angular.module('fabsmith.login', ['ngRoute', 'ngCookies'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/login', {
    templateUrl: 'ng-modules/login/login.html',
    controller: 'LoginCtrl'
  });
}])

.controller('LoginCtrl', ['$scope', '$http', '$location', function($scope, $http, $location) {
	
	// Local login function - if we do it by entering username and 
	// password in the browser
	$scope.login = function() {
		// Attempt to login via API
		$http({
			method: 'POST',
			url: '/api/users/login',
			data: {
				username: $scope.username,
				password: $scope.password,
				anticache: new Date().getTime()
			}
		})
		.success(function(data, status) {
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
					$location.path('/machines');
				})
				.error(function(data, status){
					console.log('Status: ' + status);
					console.log('Data' + data);
					alert('Could not get user data');
				});
				
			} // if data.UserId
		})
		.error(function(data, status) {
			console.log('fail code: ' + status);
			alert('Failed to log in');
		});
	};

	// Make the main controller scope accessible from outside
	// So we can use our Android app to call login function
	window.LOGIN_CTRL_SCOPE = $scope;

	// Call this from Android app as LOGIN_CTRL_SCOPE.login("user", "pass");
	$scope.androidLogin = function(username, password) {
		$scope.username = username;
		$scope.password = password;
		$scope.login();
	};

}]);

})();
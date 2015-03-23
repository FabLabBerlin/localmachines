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

				// Get user data
				$http({
					method: 'GET',
					url: '/api/users/' + data.UserId
				})
				.success(function(userData){
					$scope.onUserDataLoaded(userData);
				})
				.error(function(data, status){
					toastr.error('Could not get user data');
				});
				
			} // if data.UserId
		})
		.error(function(data, status) {
			toastr.error('Failed to log in');
		});
	};

	$scope.onUserDataLoaded = function(userData){
		// Get user roles
		$http({
			method: 'GET',
			url: '/api/users/' + userData.Id + '/roles',
			data: {
				anticache: new Date().getTime()
			}
		})
		.success(function(userRoles){
			$scope.onUserRolesLoaded(userRoles, userData);
		})
		.error(function(data, status){
			alert('Could not load user roles');
		});
	};

	$scope.onUserRolesLoaded = function(userRoles, userData){
		userData = _.merge(userData, userRoles); 
		$scope.$emit('user-login', userData);
		$location.path('/machines');
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
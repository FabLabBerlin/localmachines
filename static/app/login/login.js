(function(){

'use strict';

angular.module('fabsmith.login', ['ngRoute'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/login', {
    templateUrl: '/static/app/login/login.html',
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
			} else if (data.Status === 'logged' || data.Status == 'ok'){
				$location.path('/machines');
			}
		})
		.error(function() {
			alert('Failed to log in');
		});
	}

	// Make the main controller scope accessible from outside
	// So we can use our Android app to call login function
	window.LOGIN_CTRL_SCOPE = $scope;

	// Call this from Android app as LOGIN_CTRL_SCOPE.login("user", "pass");
	$scope.androidLogin = function(username, password) {
		$scope.username = username;
		$scope.password = password;
		$scope.login();
	}

}]);

})();
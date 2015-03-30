(function(){

'use strict';

// Declare app level module which depends on views, and components
var app = angular.module('fabsmith', [
  'ngRoute',
  'fabsmith.admin.login',
  'fabsmith.admin.dashboard',
  'fabsmith.backoffice.mainmenu',
  'fabsmith.backoffice.user',
  'fabsmith.admin.activations',
  'fabsmith.admin.machines',
  'fabsmith.admin.memberships',
  'fabsmith.admin.bookings',
  'fabsmith.admin.invoices',
  'fabsmith.admin.users',
  'fabsmith.version'
]);

// This checks whether an user is logged in always before switching to a new view
app.run(['$rootScope', '$location', '$http', function($rootScope, $location, $http) {
	$rootScope.$on('$locationChangeStart', function(event, newUrl, oldUrl) {
		
		// Get requested angular path
		var newPath = newUrl.split('#')[1];
		
		// If it is not login (main) view, 
		// check if the user is logged in
		// TODO: figure out a way that does not need a request
		if (newPath !== '/login') {
			$http({
				method: 'POST',
				url: '/api/users/login',
				params: {
					username: 'blank', // TODO: randomize?
					password: 'blank',
					anticache: new Date().getTime()
				}
			})
			.success(function(data) {
				if (data.Status !== 'logged') {
					event.preventDefault();
					$location.path('/login');
				} else {
					if (newPath) {
						$location.path(newPath);
					} else {
						$location.path('/dashboard');
					}
				}
			})
			.error(function() {
				event.preventDefault();
				$location.path('/login');
			});
		}

	});
}]); // app.run

// Configure http provider to send regular form POST data instead of JSON
app.config(['$httpProvider', function($httpProvider) {
	$httpProvider.defaults.transformRequest = function(data){
		if (data === undefined) {
		    return data;
		}
		return $.param(data);
	};

	$httpProvider.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded; charset=UTF-8';
}]); // app.config

// Main controller, checks if user logged in
app.controller('MainCtrl', ['$scope', '$http', '$location', '$cookieStore', '$cookies', 
 function($scope, $http, $location, $cookieStore, $cookies){

 	// Configure toastr default location
 	toastr.options.positionClass = 'toast-bottom-right';

	// Configure vex theme
	vex.defaultOptions.className = 'vex-theme-plain';

	// Store user data on user login
	$scope.putUserData = function(data) {
		for (var key in data) {
			if (data.hasOwnProperty(key)) {
				$cookieStore.put(key, data[key]);
			}
		}
	};
	$scope.$on('user-login', function (event, data){
		$scope.putUserData(data);
	});

	// Clear user data on user logout
	$scope.removeUserData = function() {
		for (var key in $cookies) {
			if ($cookies.hasOwnProperty(key)) {
				$cookieStore.remove(key);
			}
		}
	};

	// Log out on logout event
	$scope.logout = function() {
		$scope.removeUserData();

		$http({
			method: 'GET',
			url: '/api/users/logout',
			params: {
				anticache: new Date().getTime()
			}
		})
		.success(function() {
			$location.path('/');
		})
		.error(function() {
			alert('Failed to log out. Probably server down.');
			$location.path('/');
		});

	};
	$scope.$on('logout', function (event, data){
		$scope.logout();
	});
}]);

})(); // closure

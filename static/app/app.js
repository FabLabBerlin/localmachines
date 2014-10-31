(function(){

'use strict';

// Declare app level module which depends on views, and components
angular.module('fabsmith', [
  'ngRoute',
  'fabsmith.login',
  'fabsmith.machines',
  'fabsmith.logout',
  'fabsmith.version'
])

// Default redirect to login view
/*
.config(['$routeProvider', function($routeProvider) {
  $routeProvider.otherwise({redirectTo: '/login'});
}])
*/

// This checks whether an user is logged in always before 
// switching to a new view
.run(['$rootScope', '$location', '$http', 
	function($rootScope, $location, $http) {
	$rootScope.$on('$locationChangeStart', function(event, newUrl, oldUrl) {
		var newPath = newUrl.split('#')[1];
		if (newPath !== '/login') {
			$http({
				method: 'POST',
				url: '/api/login',
				params: {
					username: 'blank',
					password: 'blank',
					anticache: new Date().getTime()
				}
			})
			.success(function(data) {
				if (data.Status !== 'logged') {
					event.preventDefault();
					$location.path('/login');
				} else {
					// User is logged - show requested route (passthru)
				}
			})
			.error(function() {
				event.preventDefault();
				$location.path('/login');
			});
		} else {
			// Route path is /login - let it show up passthru()
		}
	});
}])

// Configure http provider to send regular form POST data instead of JSON
.config(['$httpProvider', function($httpProvider) {
	$httpProvider.defaults.transformRequest = function(data){
		if (data === undefined) {
		    return data;
		}
		return $.param(data);
	}
	$httpProvider.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded; charset=UTF-8';
}])

// Main controller, checks if user logged in
.controller('MainCtrl', ['$scope', '$http', '$location', function($scope, $http, $location){
	$http.post('/api/login')
	.success(function(data) {
		if (data.Status === 'logged') {
			$location.path('/machines');
		}
	})
	.error(function() {
		alert('Could not communicate with the API');
	});
}]);

})();

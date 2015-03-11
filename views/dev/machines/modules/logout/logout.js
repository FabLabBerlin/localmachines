'use strict';

angular.module('fabsmith.logout', ['ngRoute', 'timer'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/logout', {
    templateUrl: '/machines-view/modules/logout/logout.html',
    controller: 'LogoutCtrl'
  });
}])

.controller('LogoutCtrl', ['$scope', '$http', '$location', function($scope, $http, $location) {
	//$scope.$on('timer-stopped', function (event, data){
		//$scope.logout();
  //});

	// Activate countdown
	$scope.abortLogout = function() {
		$scope.$broadcast('timer-clear');
		$location.path('/machines');
	};

	$scope.logout = function() {
		$http({
			method: 'GET',
			url: '/api/logout',
			params: {
				anticache: new Date().getTime()
			}
		})
		.success(function() {
			$location.path('/');
		})
		.error(function() {
			alert('Failed to log out. Probably server down.');
		});
	};
}]);
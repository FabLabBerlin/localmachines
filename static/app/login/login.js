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
	$scope.date = new Date();

	$scope.login = function() {
		// Attempt to login via API
		$http.post('/api/login', {
			username: $scope.username,
			password: md5($scope.password)
		})
		.success(function(data) {
			if (data.Status === 'error') {
				alert(data.Message);
			} else if (data.Status === 'logged' || data.Status == 'ok'){
				$location.path('/machines');
			}
		})
		.error(function() {
			alert('Fatal error');
		});
	}
}]);

})();
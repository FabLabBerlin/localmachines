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
	$scope.userId = $routeParams.userId;
	console.log('$scope.userId: ', $scope.userId);
}]); // app.controller

})(); // closure
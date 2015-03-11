(function(){

'use strict';

var app = angular.module('fabsmith.backoffice.dashboard', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/dashboard', {
    templateUrl: '/admin-view/dashboard/dashboard.html',
    controller: 'DashboardCtrl'
  });
}]); // app.config

app.controller('DashboardCtrl', ['$scope', function($scope) {
	$scope.testVar = 'This is a test variable';
}]); // app.controller

})(); // closure
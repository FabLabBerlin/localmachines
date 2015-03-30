(function(){

'use strict';

var app = angular.module('fabsmith.admin.dashboard', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
 $routeProvider.when('/dashboard', {
    templateUrl: 'ng-modules/dashboard/dashboard.html',
    controller: 'DashboardCtrl'
  });
}]); // app.config

app.controller('DashboardCtrl', ['$scope', '$http', '$location', 
 function($scope, $http, $location) {

 	// Do stuff here
	
}]); // app.controller

})(); // closure
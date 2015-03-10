(function(){

'use strict';

var app = angular.module('fabsmith.backoffice.dashboard', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/dashboard', {
    templateUrl: '/static/dev/admin/dashboard/dashboard.html',
    controller: 'DashboardCtrl'
  });
}]); // app.config

})(); // closure
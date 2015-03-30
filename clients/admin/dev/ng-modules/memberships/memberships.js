(function(){

'use strict';

var app = angular.module('fabsmith.admin.memberships', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
    $routeProvider.when('/memberships', {
        templateUrl: 'ng-modules/memberships/memberships.html',
        controller: 'MembershipsCtrl'
    });
}]); // app.config

app.controller('MembershipsCtrl', ['$scope', '$http', '$location', 
 function($scope, $http, $location) {

 	// Stuff goes here

}]); // app.controller

})(); // closure
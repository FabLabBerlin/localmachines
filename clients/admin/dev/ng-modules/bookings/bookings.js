(function(){

'use strict';

var app = angular.module('fabsmith.admin.bookings', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
    $routeProvider.when('/bookings', {
        templateUrl: 'ng-modules/bookings/bookings.html',
        controller: 'BookingsCtrl'
    });
}]); // app.config

app.controller('BookingsCtrl', ['$scope', '$http', '$location', 
 function($scope, $http, $location) {

 	// Stuff goes here

}]); // app.controller

})(); // closure
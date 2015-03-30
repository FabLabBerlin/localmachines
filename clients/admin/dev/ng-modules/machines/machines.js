(function(){

'use strict';

var app = angular.module('fabsmith.admin.machines', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
    $routeProvider.when('/machines', {
        templateUrl: 'ng-modules/machines/machines.html',
        controller: 'MachinesCtrl'
    });
}]); // app.config

app.controller('MachinesCtrl', ['$scope', '$http', '$location', 
 function($scope, $http, $location) {

 	// Stuff goes here

}]); // app.controller

})(); // closure
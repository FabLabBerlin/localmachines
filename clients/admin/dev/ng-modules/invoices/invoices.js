(function(){

'use strict';

var app = angular.module('fabsmith.admin.invoices', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
    $routeProvider.when('/invoices', {
        templateUrl: 'ng-modules/invoices/invoices.html',
        controller: 'InvoicesCtrl'
    });
}]); // app.config

app.controller('InvoicesCtrl', ['$scope', '$http', '$location', 
 function($scope, $http, $location) {

 	// Stuff goes here

}]); // app.controller

})(); // closure
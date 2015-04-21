(function(){

'use strict';

var app = angular.module('fabsmith.admin.invoices', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
    $routeProvider.when('/invoices', {
        templateUrl: 'ng-modules/invoices/invoices.html',
        controller: 'InvoicesCtrl'
    });
}]); // app.config

app.filter('myDate', function(){
  return function(val){
    var date = new Date(val);
    return date;
  };    
});

app.controller('InvoicesCtrl', ['$scope', '$http', '$location', 
 function($scope, $http, $location) {

 	$scope.invoices = [];

 	$scope.offset = new Date().getTimezoneOffset();

 	// Load invoices
 	$http({
 		method: 'GET',
 		url: '/api/invoices'
 	})
 	.success(function(invoices) {
 		toastr.success('Invoices loaded');
 		$scope.invoices = invoices;
 	})
 	.error(function() {
 		toastr.error('Failed to load invoices');
 	});

}]); // app.controller

})(); // closure
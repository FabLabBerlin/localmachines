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

  // Load invoices
  $scope.loadInvoices = function() {
    $http({
      method: 'GET',
      url: '/api/invoices'
    })
    .success(function(invoices) {
      $scope.invoices = invoices;
    })
    .error(function() {
      toastr.error('Failed to load invoices');
    });
  };

  $scope.invoices = [];
  $scope.loadInvoices();

  $scope.deleteInvoice = function(invoiceId) {
    $http({
      method: 'DELETE',
      url: '/api/invoices/' + invoiceId
    })
    .success(function(response) {
      $scope.loadInvoices();
    })
    .error(function() {
      toastr.error("Failed to delete invoice");
    });
  };

}]); // app.controller

})(); // closure
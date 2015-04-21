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

app.controller('InvoicesCtrl', ['$scope', '$http', '$location', 'randomToken', 
 function($scope, $http, $location, randomToken) {

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

  $scope.deleteInvoicePrompt = function(invoiceId) {
    var token = randomToken.generate();
    vex.dialog.prompt({
      message: 'Enter <span class="delete-prompt-token">' + 
       token + '</span> to delete',
      placeholder: 'Token',
      callback: $scope.deleteInvoicePromptCallback.bind(this, token, invoiceId)
    });
  };

  $scope.deleteInvoicePromptCallback = function(expectedToken, invoiceId, value) {
    if (value) {    
      if (value === expectedToken) {
        $scope.deleteInvoice(invoiceId);
      } else {
        toastr.error('Wrong token');
      }
    } else if (value !== false) {
      toastr.error('No token');
    }
  };

  $scope.deleteInvoice = function(invoiceId) {
    $http({
      method: 'DELETE',
      url: '/api/invoices/' + invoiceId
    })
    .success(function(response) {
      toastr.success("Invoice deleted");
      $scope.loadInvoices();
    })
    .error(function() {
      toastr.error("Failed to delete invoice");
    });
  };

}]); // app.controller

})(); // closure
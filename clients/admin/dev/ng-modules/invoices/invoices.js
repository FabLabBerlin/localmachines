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

app.controller('InvoicesCtrl', ['$scope', '$http', '$location', '$cookies', 'randomToken',
 function($scope, $http, $location, $cookies, randomToken) {

  // Load invoices
  $scope.loadInvoices = function() {
    $http({
      method: 'GET',
      url: '/api/invoices',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(invoices) {
      $scope.invoices = invoices;
      _.each(invoices, function(invoice) {
        $scope.invoicesById[invoice.Id] = invoice;
      });
    })
    .error(function() {
      toastr.error('Failed to load invoices');
    });
  };

  $scope.invoices = [];
  $scope.invoicesById = {};
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
      url: '/api/invoices/' + invoiceId,
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(response) {
      toastr.success("Invoice deleted");
      $scope.loadInvoices();
    })
    .error(function() {
      toastr.error("Failed to delete invoice");
    });
  };

  $scope.createFbDraftsPrompt = function(invoiceId) {
    var token = randomToken.generate();
    vex.dialog.prompt({
      message: 'Enter <span class="delete-prompt-token">' +
      token + '</span> to create Fastbill drafts',
      placeholder: 'Token',
      callback: function(value) {
        if (value) {
          if (value === token) {
            $scope.createFbDrafts(invoiceId);
          } else {
            toastr.error('Wrong token');
          }
        } else if (value !== false) {
          toastr.error('No token');
        }
      } // callback
    });
  };

  $scope.createFbDrafts = function(invoiceId) {
    var params = {
      startDate: $scope.invoicesById[invoiceId].PeriodFrom,
      endDate: $scope.invoicesById[invoiceId].PeriodTo,
      location: $cookies.locationId,
      ac: new Date().getTime()
    };
    $http({
      method: 'POST',
      url: '/api/invoices/create_drafts',
      params: params
    })
    .success(function() {
      toastr.info('Sucessfully created invoice drafts');
    })
    .error(function() {
      toastr.error('Error creating invoice');
    });
  };

}]); // app.controller

})(); // closure
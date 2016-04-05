(function(){

'use strict';

var app = angular.module('fabsmith.admin.invoices', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/invoices', {
    templateUrl: 'ng-modules/invoices/invoices.html',
    controller: 'InvoicesCtrl'
  });
}]); // app.config

app.controller('InvoicesCtrl',
 ['$scope', '$http', '$location', '$cookies', 'api', 'randomToken', 
 function($scope, $http, $location, $cookies, api, randomToken) {

  $scope.invoices = [];
  $scope.loading = false;

  // Load invoices
  $scope.loadInvoices = function() {
    $http({
      method: 'GET',
      url: '/api/invoices',
      params: {
        ac: new Date().getTime(),
        location: $cookies.locationId
      }
    })
    .success(function(invoices) {
      $scope.invoices = invoices;
    })
    .error(function() {
      toastr.error('Failed to load invoices');
    });
  };

  $scope.createFbDraftsPrompt = function(id) {
    var token = randomToken.generate();
    vex.dialog.prompt({
      message: 'Enter <span class="delete-prompt-token">' +
      token + '</span> to create Fastbill drafts',
      placeholder: 'Token',
      callback: function(value) {
        if (value) {
          if (value === token) {
            $scope.createFbDrafts(id);
          } else {
            toastr.error('Wrong token');
          }
        } else if (value !== false) {
          toastr.error('No token');
        }
      } // callback
    });
  };

  $scope.createFbDrafts = function(id) {
    console.log('invoices: ', $scope.invoices);
    $scope.loading = true;
    $http({
      method: 'POST',
      url: '/api/invoices/' + id + '/create_drafts',
      params: {
        location: $cookies.locationId
      }
    })
    .success(function(draftsReport) {
      $scope.loading = false;
      console.log('draftsReport=', draftsReport);
      console.log('$scope.usersById=', $scope.usersById);
      $scope.draftsReport = draftsReport;
      toastr.info('Sucessfully created invoice drafts');
    })
    .error(function() {
      $scope.loading = false;
      toastr.error('Error creating invoice');
    });
  };

  api.loadUsers(function(userData) {
    $scope.users = userData.users;
    $scope.usersById = userData.usersById;
    $scope.loadInvoices();
  });

}]); // app.controller

})(); // closure
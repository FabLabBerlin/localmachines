(function(){

'use strict';

var app = angular.module('fabsmith.admin.activations', 
 ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/activations', {
    templateUrl: 'ng-modules/activations/activations.html',
    controller: 'ActivationsCtrl'
  });
}]); // app.config

app.controller('ActivationsCtrl',
 ['$scope', '$http', '$location', '$cookies', 'randomToken', 'api',
 function($scope, $http, $location, $cookies, randomToken, api) {

  $scope.activationsStartDate = moment().format('YYYY-MM');
  $scope.activationsEndDate = moment().format('YYYY-MM');
  $scope.searchTeam = '';
  $scope.usersById = {};
  $scope.loading = false;

  $scope.loadPage = function() {
    var offset = $scope.itemsPerPage * ($scope.currentPage - 1);
    $scope.pageActivations = $scope.activations.slice(offset, offset + $scope.itemsPerPage);
  };

  // Loads and reloads activations according to filter.
  // If user ID is not set - load all users
  $scope.loadActivations = function() {
    
    if (!$scope.machines) {
      toastr.error('Machines not loaded');
      return;
    }

    if ($scope.activationsStartDate.length !== 7 || $scope.activationsEndDate.length !== 7) {
      return;
    }

    $http({
      method: 'GET',
      url: '/api/activations', 
      params: {
        startDate: $scope.activationsStartDate,
        endDate: $scope.activationsEndDate,
        search: $scope.searchTerm,
        userId: 1,
        includeInvoiced: false,
        itemsPerPage: $scope.itemsPerPage,
        page: $scope.currentPage,
        ac: new Date().getTime(),
        location: $cookies.get('location')
      }
    })
    .success(function(activations) {
      _.each(activations, function(activation) {
        var machine = _.find($scope.machines, 'Id', activation.MachineId);
        if (machine) {
          activation.MachineName = machine.Name;
        } else {
          activation.MachineName = 'Undefined';
        }
      });

      loadUserNames();

      $scope.activations = activations;
      $scope.numActivations = activations.length;
      $scope.numPages = Math.ceil($scope.numActivations / $scope.itemsPerPage);
      $scope.loadPage();
    })
    .error(function() {
      toastr.error('Failed to load activations');
    });
  };

  function loadUserNames(userId) {
    $http({
      method: 'GET',
      url: '/api/users',
      params: {
        ac: new Date().getTime(),
        location: $cookies.get('location')
      }
    })
    .success(function(users) {
      var usersById = {};
      _.each(users, function(user) {
        $scope.usersById[user.Id] = user;
      });
      _.each($scope.activations, function(activation) {
        var user = $scope.usersById[activation.UserId];
        if (user) {
          activation.UserName = user.FirstName + ' ' + user.LastName;
        }
      });
    })
    .error(function() {
      toastr.error('Failed to load user name');
    });
  }

  // This is called whenever start or end date changes
  $scope.onFilterChange = function() {
    if ($scope.activationsStartDate && $scope.activationsEndDate) {
      $scope.activations = [];
      $scope.currentPage = 1;
      $scope.loadActivations();
    }
  };

  $scope.activations = [];
  $scope.currentPage = 1;
  $scope.itemsPerPage = 10;

  $scope.loadNextPage = function() {
    if ($scope.activations.length < $scope.itemsPerPage) {
      return;
    }

    $scope.currentPage++;
    $scope.loadPage();
  };

  $scope.loadPrevPage = function() {
    if ($scope.currentPage < 1) {
      return;
    }

    $scope.currentPage--;
    $scope.loadPage();
  };

  $scope.createActivation = function() {
    $http({
      method: 'POST',
      url: '/api/activations',
      params: {
        location: $cookies.get('location')
      }
    })
    .success(function(a) {
      $location.path('/activations/' + a.Id);
    })
    .error(function() {
      toastr.error('Error.  Please try again later.');
    });
  };

  function getExportParams() {
    return {
      startDate: $scope.activationsStartDate,
      endDate: $scope.activationsEndDate,
      //userId: 1,
      includeInvoiced: false,
      itemsPerPage: $scope.itemsPerPage,
      page: $scope.currentPage,
      location: $cookies.get('location'),
      ac: new Date().getTime()
    };
  }

  // Creates invoice on the server side and returns link
  $scope.exportSpreadsheet = function() {
    $http({
      method: 'POST',
      url: '/api/billing/monthly_earnings',
      params: getExportParams()
    })
    .success(function(invoiceData) {
      // invoiceData should contain link to the generated xls file
      toastr.success('Invoice created');
      console.log(invoiceData);

      var filePathParts = invoiceData.FilePath.split("/");
      var fileName = filePathParts[filePathParts.length-1];

      var alertContent = '<div class="row">'+
        '<div class="col-xs-6"><b>Invoice created!</b><br>'+
        '<b>File name:</b> ' + fileName + '</div>'+
        '<div class="col-xs-6"><a '+
        'href="/api/invoices/' + invoiceData.Id + '/download_excel" '+
        'class="btn btn-primary btn-block">'+
        'Download</a></div>'+
        '</div>';
      vex.dialog.alert(alertContent);
    })
    .error(function() {
      toastr.error('Error creating invoice');
    });
  };

  $scope.edit = function(id) {
    $location.path('/activations/' + id);
  };

  // We need full machine names for the activation table
  if (!$scope.machines) {
    api.loadMachines(function(resp) {
      $scope.machines = resp.machines;
      $scope.loadActivations();
    });
  }

}]); // app.controller

})(); // closure
(function(){

'use strict';
var app = angular.module('fabsmith.admin.coworking.rental', ['ngRoute', 'ngCookies', 'fabsmith.admin.randomtoken', 'fabsmith.admin.api']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/coworking/rentals/:id', {
    templateUrl: 'ng-modules/coworking/rental.html',
    controller: 'RentalCtrl'
  });
}]); // app.config

app.controller('RentalCtrl',
 ['$scope', '$routeParams', '$http', '$location', 'randomToken', 'api',
 function($scope, $routeParams, $http, $location, randomToken, api) {

  $scope.purchases = [];
  $scope.rental = {
    Id: $routeParams.id
  };
  $scope.tablesById = {};
  $scope.users = [];
  $scope.usersById = {};

  function loadTables() {
    $http({
      method: 'GET',
      url: '/api/products',
      params: {
        ac: new Date().getTime(),
        type: 'co-working'
      }
    })
    .success(function(data) {
      $scope.tables = _.sortBy(data, function(table) {
        return table.Product.Name;
      });
      _.each($scope.tables, function(table) {
        $scope.tablesById[table.Product.Id] = table;
      });
      loadTablePurchase();
    })
    .error(function() {
      toastr.error('Failed to get tables');
    });
  }

  function loadTablePurchase() {
    $http({
      method: 'GET',
      url: '/api/purchases/' + $scope.rental.Id,
      params: {
        ac: new Date().getTime(),
        type: 'co-working'
      }
    })
    .success(function(sp) {
      $scope.rental = sp;
      var start = moment(sp.TimeStart).tz('Europe/Berlin');
      var end = moment(sp.TimeEnd).tz('Europe/Berlin');
      sp.DateStartLocal = start.format('YYYY-MM-DD');
      sp.DateEndLocal = end.format('YYYY-MM-DD');
      sp.TimeStartLocal = start.format('HH:mm');
      sp.TimeEndLocal = end.format('HH:mm');
      calculateTotalPrice();
      setTimeout(function() {
        $('.selectpicker').selectpicker('refresh');
      }, 100);
    })
    .error(function(data, status) {
      toastr.error('Failed to load user data');
    });
  }

  function parseInputTimes() {
    var sp = $scope.rental;
    sp.TimeStart = moment.tz(sp.DateStartLocal + ' ' + sp.TimeStartLocal, 'Europe/Berlin').toDate();
    sp.TimeEnd = moment.tz(sp.DateEndLocal + ' ' + sp.TimeEndLocal, 'Europe/Berlin').toDate();
  }

  function calculateQuantity() {
    console.log('$scope.timeChange()');
    parseInputTimes();
    var start = moment($scope.rental.TimeStart);
    var end = moment($scope.rental.TimeEnd);
    var duration = end.unix() - start.unix();
    console.log('duration=', duration);
    var quantity;
    switch ($scope.rental.PriceUnit) {
    case 'minute':
      quantity = duration / 60;
      break;
    case 'hour':
      quantity = duration / 3600;
      break;
    case 'day':
      quantity = duration / 24 / 3600;
      break;
    default:
      return;
    }
    $scope.rental.Quantity = quantity;
  }
   
  // https://www.artstation.com/artwork/b5zBn
  function calculateTotalPrice() {
    var totalPrice = $scope.rental.Quantity * $scope.rental.PricePerUnit;
    $scope.rental.TotalPrice = totalPrice.toFixed(2);
  }

  $scope.tableChange = function() {
    $('#rental-table').selectpicker('refresh');
    console.log('$scope.tableChange()');
    var tableId = parseInt($scope.rental.ProductId);
    var table = $scope.tablesById[tableId];
    $scope.rental.PricePerUnit = table.Product.Price;
    $scope.rental.PriceUnit = table.Product.PriceUnit;
    calculateQuantity();
    calculateTotalPrice();    
  };

  $scope.userChange = function() {
    $('#rental-user').selectpicker('refresh');
  };

  $scope.timeChange = function() {
    calculateQuantity();
    calculateTotalPrice();
  };

  $scope.priceUnitChange = function() {
    calculateQuantity();
    calculateTotalPrice();
  };

  $scope.save = function() {
    parseInputTimes();
    $http({
      method: 'PUT',
      url: '/api/purchases/' + $scope.rental.Id + '?type=co-working',
      headers: {'Content-Type': 'application/json' },
      data: $scope.rental,
      transformRequest: function(data) {
        var transformed = _.extend({}, data);
        transformed.ProductId = parseInt(data.ProductId);
        transformed.UserId = parseInt(data.UserId);
        transformed.Quantity = parseInt(data.Quantity);
        transformed.PricePerUnit = parseInt(data.PricePerUnit);
        transformed.TotalPrice = parseFloat(data.TotalPrice);
        console.log('transformed:', transformed);
        return JSON.stringify(transformed);
      },
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function() {
      $location.path('/coworking');
    })
    .error(function(data) {
      toastr.error('Error while trying to save changes');
    });
  };

  // Init scope variables
  var pickadateOptions = {
    format: 'yyyy-mm-dd'
  };
  $('.datepicker').pickadate(pickadateOptions);

  api.loadUsers(function(userData) {
    $scope.users = userData.users;
    $scope.usersById = userData.usersById;
    loadTables();
  });

}]); // app.controller

})(); // closure

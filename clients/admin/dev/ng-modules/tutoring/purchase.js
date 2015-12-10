(function(){

'use strict';

var app = angular.module('fabsmith.admin.tutoring.purchase',
  ['ngRoute', 'ngCookies', 'fabsmith.admin.api']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/tutoring/purchases/:id', {
    templateUrl: 'ng-modules/tutoring/purchase.html',
    controller: 'PurchaseCtrl'
  });
}]); // app.config

app.controller('PurchaseCtrl', ['$scope', '$routeParams', '$http', '$location', '$filter', 'api',
  function($scope, $routeParams, $http, $location, $filter, api) {

  $scope.purchase = {
    Id: $routeParams.id
  };

  api.loadTutors(function(tutorData) {
    $scope.tutors = tutorData.tutors;
    $scope.tutorsById = tutorData.tutorsById;
    console.log('$scope.tutors:', $scope.tutors);
    api.loadUsers(function(userData) {
      $scope.users = userData.users;
      $scope.usersById = userData.usersById;
      api.loadTutoringPurchase($scope.purchase.Id, function(purchase) {
        $scope.purchase = purchase;
        console.log('$scope.purchase = ', $scope.purchase);
        setTimeout(function() {
          $('.selectpicker').selectpicker('refresh');
        }, 100);
        //api.purchase.calculateQuantity($scope.purchase);
        // Why would I have to search for this in another place?
        //api.purchase.calculateTotalPrice($scope.purchase);

        calculateDurations();

        // We do not include the following in the calculateDurations
        // because we won't be able to change the values in the input field.
        var start = moment($scope.purchase.TimeStart);
        $scope.purchase.DateStartLocal = start.format('YYYY-MM-DD');
        $scope.purchase.TimeStartLocal = start.format('HH:mm');

        var end = moment($scope.purchase.TimeEnd);
        $scope.purchase.DateEndLocal = end.format('YYYY-MM-DD');
        $scope.purchase.TimeEndLocal = end.format('HH:mm');
      });
    });
  });

  // Init scope variables
  var pickadateOptions = {
    format: 'yyyy-mm-dd'
  };
  $('.datepicker').pickadate(pickadateOptions);

  $scope.tutorChange = function() {
    $('#tp-tutor').selectpicker('refresh');
    var tutorId = parseInt($scope.purchase.ProductId);
    var tutor = $scope.tutorsById[tutorId];
    console.log('tutor:', tutor);
    $scope.purchase.PricePerUnit = tutor.Product.Price;
    $scope.purchase.PriceUnit = tutor.Product.PriceUnit;
    calculateDurations();
  };

  $scope.userChange = function() {
    $('#tp-user').selectpicker('refresh');
  };

  $scope.timeChange = function() {
    calculateDurations();
  };

  function calculateDurations() {
    console.log('calculateDurations()');

    // These are tutoring specific things,
    // why would one need to search for them in another file/place?
    //api.purchase.calculateQuantity($scope.purchase);
    //api.purchase.calculateTotalPrice($scope.purchase);

    // Combine the date and time of start and end time
    var startTimeString = $scope.purchase.DateStartLocal + ' ' + $scope.purchase.TimeStartLocal;
    var endTimeString = $scope.purchase.DateEndLocal + ' ' + $scope.purchase.TimeEndLocal;

    var start = moment(startTimeString);
    var end = moment(endTimeString);
    
    if (start.unix() > 0) {
      if (end.unix() > 0) {
        var reservedDuration = moment.duration(end.diff(start));
        $scope.purchase.TimeReserved = reservedDuration.format('h[h] m[m]');
      }
      if ($scope.purchase.Quantity > 0) {
        //$scope.purchase.TimeTimed = end.clone().subtract(start).format('HH:mm:ss');
        var momentDuration;
        
        if ($scope.purchase.PriceUnit === 'hour') {
          momentDuration = moment.duration($scope.purchase.Quantity, 'hours');
        } else if ($scope.purchase.PriceUnit === 'minute') {
          momentDuration = moment.duration($scope.purchase.Quantity, 'minutes');
        } else if ($scope.purchase.PriceUnit === 'day') {
          momentDuration = moment.duration($scope.purchase.Quantity, 'days');
        }

        if (momentDuration.asHours() < 1) {
          $scope.purchase.TimeTimed = momentDuration.format('m[m] s[s]');
        } else {
          $scope.purchase.TimeTimed = momentDuration.format('h[h] m[m]');
        }

        $scope.purchase.TotalPrice = $scope.purchase.Quantity * $scope.purchase.PricePerUnit;
        $scope.purchase.TotalPrice = $filter('currency')($scope.purchase.TotalPrice, '');
      }
    }
  }

  $scope.priceUnitChange = function() {
    api.purchase.calculateQuantity($scope.purchase);
    api.purchase.calculateTotalPrice($scope.purchase);
  };


  $scope.save = function() {
    api.purchase.parseInputTimes($scope.purchase);
    $http({
      method: 'PUT',
      url: '/api/purchases/' + $scope.purchase.Id + '?type=tutor',
      headers: {'Content-Type': 'application/json' },
      data: $scope.purchase,
      transformRequest: function(data) {
        var transformed = _.extend({}, data);
        transformed.ProductId = parseInt(data.ProductId);
        transformed.UserId = parseInt(data.UserId);
        transformed.Quantity = parseFloat(data.Quantity);
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
      $location.path('/tutoring');
    })
    .error(function(data) {
      toastr.error('Error while trying to save changes');
    });
  };

  $scope.cancel = function() {
    $location.path('/tutoring');
  };

}]); // app.controller

})(); // closure
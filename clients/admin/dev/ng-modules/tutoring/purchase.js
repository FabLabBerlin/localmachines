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

app.controller('PurchaseCtrl', ['$scope', '$routeParams', '$http', '$location', 'api',
  function($scope, $routeParams, $http, $location, api) {

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
        api.purchase.calculateQuantity($scope.purchase);
        api.purchase.calculateTotalPrice($scope.purchase);
        calculateDurations();
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
    api.purchase.calculateQuantity($scope.purchase);
    api.purchase.calculateTotalPrice($scope.purchase);
    var start = moment($scope.purchase.TimeStart);
    var end = moment($scope.purchase.TimeEnd);
    var endPlanned = moment($scope.purchase.TimeEndPlanned);
    console.log('start:', start.format('YYYY-MM-DD HH:mm'));
    console.log('end:', end.format('YYYY-MM-DD HH:mm'));
    console.log('endPlanned:', endPlanned.format('YYYY-MM-DD HH:mm'));
    if (start.unix() > 0) {
      if (endPlanned.unix() > 0) {
        console.log('setting TimeReserved...');
        $scope.purchase.TimeReserved = endPlanned.clone().subtract(start).format('HH:mm:ss');
        console.log('$scope.purchase.TimeReserved:', $scope.purchase.TimeReserved);
      }
      if (end.unix() > 0) {
        console.log('setting TimeTimed...');
        $scope.purchase.TimeTimed = end.clone().subtract(start).format('HH:mm:ss');
        console.log('$scope.purchase.TimeTimed:', $scope.purchase.TimeTimed);
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
(function(){

'use strict';

var app = angular.module('fabsmith.admin.reservations', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/reservations', {
    templateUrl: 'ng-modules/reservations/reservations.html',
    controller: 'ReservationsCtrl'
  });
}]); // app.config

app.controller('ReservationsCtrl', ['$scope', '$http', '$location', '$cookieStore', 
 function($scope, $http, $location, $cookieStore) {

  $scope.machines = [];
  $scope.reservationRules = [];
  $scope.machinesById = {};

  function loadReservationRules() {
    $http({
      method: 'GET',
      url: '/api/reservation_rules',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      $scope.reservationRules = _.map(data, function(r) {
        r.Machine = $scope.machinesById[r.MachineId] || {};
        return r;
      });
    })
    .error(function() {
      toastr.error('Failed to get reservations');
    });
  }

  $http({
    method: 'GET',
    url: '/api/machines',
    params: {
      ac: new Date().getTime()
    }
  })
  .success(function(data) {
    $scope.machines = data;
    _.each($scope.machines, function(machine) {
      $scope.machinesById[machine.Id] = machine;
    });
    loadReservationRules();
  })
  .error(function() {
    toastr.error('Failed to get machines');
  });

  $scope.addReservationRule = function() {
    $http({
      method: 'POST',
      url: '/api/reservation_rules',
      headers: {'Content-Type': 'application/json' },
      data: {
        MachineId: parseInt($('select').val()),
        Created: new Date()
      },
      transformRequest: function(data) {
        var transformed = _.extend({}, data);
        return JSON.stringify(transformed);
      }
    })
    .success(function(data) {
      toastr.info('Congratulations, you have successfully created a new Reservation Rule!');
      loadReservationRules();
    })
    .error(function() {
      toastr.error('Failed to create Reservation Rule');
    });
  };

}]); // app.controller

})(); // closure
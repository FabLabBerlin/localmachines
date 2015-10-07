(function(){

'use strict';

var app = angular.module('fabsmith.admin.reservations', ['ngRoute', 'ngCookies', 'fabsmith.admin.randomtoken']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/reservations', {
    templateUrl: 'ng-modules/reservations/reservations.html',
    controller: 'ReservationsCtrl'
  });
}]); // app.config

app.controller('ReservationsCtrl', ['$scope', '$http', '$location', '$cookieStore', 'randomToken', 
 function($scope, $http, $location, $cookieStore, randomToken) {

  $scope.machines = [];
  $scope.reservationRules = [];
  $scope.machinesById = {};
  $scope.reservationRulesById = {};

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
        $scope.reservationRulesById[r.Id] = r;
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
        Created: new Date(),
        Unavailable: true
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

  $scope.saveReservationRule = function(id) {
    $http({
      method: 'PUT',
      url: '/api/reservation_rules/' + id,
      headers: {'Content-Type': 'application/json' },
      data: $scope.reservationRulesById[id],
      transformRequest: function(data) {
        var transformed = _.extend({}, data);
        return JSON.stringify(transformed);
      }
    })
    .success(function(data) {
      toastr.info('Saved updates to Reservation Rule');
    })
    .error(function() {
      toastr.error('Failed to update Reservation Rule');
    });
  };

  $scope.deleteReservationRule = function(id) {
    var token = randomToken.generate();
    vex.dialog.prompt({
      message: 'Enter <span class="delete-prompt-token">' +
       token + '</span> to delete',
      placeholder: 'Token',
      callback: function(value) {
        if (value) {
          if (value === token) {
            $http({
              method: 'DELETE',
              url: '/api/reservation_rules/' + id,
              params: {
                ac: new Date().getTime()
              }
            })
            .success(function() {
              toastr.success('Reservation Rule deleted');
              loadReservationRules();
            })
            .error(function() {
              toastr.error('Error while trying to delete Reservation Rule');
            });
          } else {
            toastr.error('Wrong token');
          }
        } else {
          toastr.warning('Deletion canceled');
        }
      }
    });
  };

  $scope.setAvailable = function(id) {
    var rule = $scope.reservationRulesById[id];
    rule.Available = true;
    rule.Unavailable = false;
  };

  $scope.setUnavailable = function(id) {
    var rule = $scope.reservationRulesById[id];
    rule.Available = false;
    rule.Unavailable = true;
  };

}]); // app.controller

})(); // closure
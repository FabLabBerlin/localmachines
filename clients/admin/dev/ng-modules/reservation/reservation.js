(function(){

'use strict';
var app = angular.module('fabsmith.admin.reservation', ['ngRoute', 'ngCookies', 'fabsmith.admin.randomtoken']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/reservations/:reservationId', {
    templateUrl: 'ng-modules/reservation/reservation.html',
    controller: 'ReservationCtrl'
  });
}]); // app.config

app.controller('ReservationCtrl',
 ['$scope', '$routeParams', '$http', '$location', 'randomToken',
 function($scope, $routeParams, $http, $location, randomToken) {

  $scope.reservation = {
    Id: $routeParams.reservationId
  };

  $scope.machines = [];
  $scope.machinesById = {};
  $scope.users = [];
  $scope.usersById = {};

  function loadMachines() {
    $http({
      method: 'GET',
      url: '/api/machines',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      $scope.machines = _.sortBy(data, function(machine) {
        return machine.Name;
      });
      _.each($scope.machines, function(machine) {
        $scope.machinesById[machine.Id] = machine;
      });
      loadReservation($scope.reservation.Id);
    })
    .error(function() {
      toastr.error('Failed to get machines');
    });
  }

  function loadReservation(id) {
    $http({
      method: 'GET',
      url: '/api/reservations/' + id,
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(r) {
      $scope.reservation = r;
      r.TimeStartLocal = moment(r.TimeStart).tz('Europe/Berlin').format('YYYY-MM-DD HH:mm');
      r.TimeEndLocal = moment(r.TimeEnd).tz('Europe/Berlin').format('YYYY-MM-DD HH:mm');
    })
    .error(function(data, status) {
      toastr.error('Failed to load user data');
    });
  }

  function loadUsers() {
    $http({
      method: 'GET',
      url: '/api/users',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      $scope.users = _.sortBy(data, function(user) {
        return user.FirstName + ' ' + user.LastName;
      });
      _.each($scope.users, function(user) {
        $scope.usersById[user.Id] = user;
      });
      loadMachines();
    })
    .error(function() {
      toastr.error('Failed to get reservations');
    });
  }

  $scope.deleteReservation = function() {
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
              url: '/api/reservations/' + $scope.reservation.Id,
              params: {
                ac: new Date().getTime()
              }
            })
            .success(function() {
              toastr.success('Reservation deleted');
              $location.path('/reservations');
            })
            .error(function() {
              toastr.error('Error while trying to delete Reservation');
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

  $scope.saveReservation = function() {
    $http({
      method: 'PUT',
      url: '/api/reservations/' + $scope.reservation.Id,
      headers: {'Content-Type': 'application/json' },
      data: $scope.reservation,
      transformRequest: function(data) {
        var transformed = _.extend({}, data);
        transformed.MachineId = parseInt(data.MachineId);
        transformed.TimeStart = moment.tz(data.TimeStartLocal, 'Europe/Berlin').toDate();
        transformed.TimeEnd = moment.tz(data.TimeEndLocal, 'Europe/Berlin').toDate();
        transformed.UserId = parseInt(data.UserId);
        console.log('transformed:', transformed);
        return JSON.stringify(transformed);
      },
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function() {
      toastr.success('Reservation updated');
    })
    .error(function(data) {
      if (data === 'duplicateEntry') {
        toastr.error('Duplicate entry error. Make sure that fields like user name and email are unique.');
      } else if (data === 'lastAdmin') {
        $scope.user.Admin = true;
        $scope.updateAdminStatus();
        toastr.error('You are the last remaining admin. Remember - power comes with great responsibility!');
      } else if (data === 'selfAdmin') {
        $scope.user.Admin = true;
        $scope.updateAdminStatus();
        toastr.error('You can not unadmin yourself. Someone else has to do it.');
      } else {
        toastr.error('Error while trying to save changes');
      }
    });
  };

  loadUsers();

}]); // app.controller

})(); // closure

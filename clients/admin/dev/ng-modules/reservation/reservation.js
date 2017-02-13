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
 ['$scope', '$routeParams', '$http', '$location', '$cookies', 'randomToken', 'api',
 function($scope, $routeParams, $http, $location, $cookies, randomToken, api) {

  $scope.reservation = {
    Id: $routeParams.reservationId
  };

  $scope.machines = [];
  $scope.machinesById = {};
  $scope.users = [];
  $scope.usersById = {};

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
      setTimeout(function() {
        $('.selectpicker').selectpicker('refresh');
      }, 100);
    })
    .error(function(data, status) {
      toastr.error('Failed to load user data');
    });
  }

  function loadMachines() {
    api.loadMachines(function(resp) {
      $scope.machines = _.sortBy(resp.machines, function(machine) {
        return machine.Name;
      });
      _.each($scope.machines, function(machine) {
        $scope.machinesById[machine.Id] = machine;
      });
      loadReservation($scope.reservation.Id);
    });
  }

  function loadUsers() {
    $http({
      method: 'GET',
      url: '/api/users',
      params: {
        ac: new Date().getTime(),
        location: $cookies.get('location')
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

  $scope.refreshUserPicker = function() {
    $('#reservation-user').selectpicker('refresh');
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

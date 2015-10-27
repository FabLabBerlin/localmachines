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
  $scope.users = [];
  $scope.usersById = {};

  /*
   *
   * Loader functions
   *
   */

  function loadUsers() {
    $http({
      method: 'GET',
      url: '/api/users',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      $scope.users = data;
      _.each($scope.users, function(user) {
        $scope.usersById[user.Id] = user;
      });
      _.each($scope.reservations, function(r) {
        r.User = $scope.usersById[r.UserId];
      });
    })
    .error(function() {
      toastr.error('Failed to get reservations');
    });
  }

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

  function loadReservations() {
    $http({
      method: 'GET',
      url: '/api/reservations',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      var reservations = _.map(data, function(r) {
        r.Machine = $scope.machinesById[r.MachineId] || {};
        r.TimeStartLocal = moment(r.TimeStart).tz('Europe/Berlin').format('YYYY-MM-DD HH:mm');
        r.TimeEndLocal = moment(r.TimeEnd).tz('Europe/Berlin').format('YYYY-MM-DD HH:mm');
        return r;
      });
      $scope.reservations = _.sortBy(reservations, function(r) {
        return [r.Machine.Name, r.TimeStart];
      });
      loadUsers();
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
      if (_.isNumber(machine.ReservationPriceHourly)) {
        machine.ReservationPriceHalfHourly = machine.ReservationPriceHourly / 2;
      }
      $scope.machinesById[machine.Id] = machine;
    });
    loadReservationRules();
    loadReservations();
  })
  .error(function() {
    toastr.error('Failed to get machines');
  });


  /*
   *
   * Machines Rules CRUD functions
   *
   */

  $scope.saveMachine = function(id) {
    var machine = _.clone($scope.machinesById[id]);

    if (machine.ReservationPriceStart) {
      machine.ReservationPriceStart = parseFloat(machine.ReservationPriceStart);
    } else {
      machine.ReservationPriceStart = null;
    }
    if (machine.ReservationPriceHalfHourly) {
      machine.ReservationPriceHourly = parseFloat(machine.ReservationPriceHalfHourly) * 2;
    } else {
      machine.ReservationPriceHourly = null;
    }

    $http({
      method: 'PUT',
      url: '/api/machines/' + id,
      headers: {'Content-Type': 'application/json' },
      data: machine,
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


  /*
   *
   * Reservation Rules CRUD functions
   *
   */

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
    if (!$scope.weekDayCheck(id)) {
      toastr.error('Failed to update Reservation Rule');
      return;
    }

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


  /*
   *
   * Reservations CRUD functions
   *
   */

  $scope.deleteReservation = function(id) {
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
              url: '/api/reservations/' + id,
              params: {
                ac: new Date().getTime()
              }
            })
            .success(function() {
              toastr.success('Reservation deleted');
              loadReservations();
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


  /*
   *
   * Reservations helper functions
   *
   */

  $scope.dateRangeUpdate = function(id) {
    var rule = $scope.reservationRulesById[id];
    var dateStart = String(rule.DateStart);
    var dateEnd = String(rule.DateEnd);
    if (dateStart.length === 10 && dateEnd.length === 10) {
      dateStart = moment(dateStart);
      dateEnd = moment(dateEnd);
      var i = 0;
      if (!dateStart.isValid() || !dateEnd.isValid()) {
        return;
      }
      if (dateStart.unix() > dateEnd.unix()) {
        toastr.error('Date Start is after Date End');
        return;
      }
      var t = dateStart.clone();
      var tString;
      var dateEndString = dateEnd.format('YYYY-MM-DD');
      rule.Monday = false;
      rule.Tuesday = false;
      rule.Wednesday = false;
      rule.Thursday = false;
      rule.Friday = false;
      rule.Saturday = false;
      rule.Sunday = false;
      do {
        tString = t.format('YYYY-MM-DD');
        switch (t.isoWeekday()) {
        case 1:
          rule.Monday = true;
          break;
        case 2:
          rule.Tuesday = true;
          break;
        case 3:
          rule.Wednesday = true;
          break;
        case 4:
          rule.Thursday = true;
          break;
        case 5:
          rule.Friday = true;
          break;
        case 6:
          rule.Saturday = true;
          break;
        case 7:
          rule.Sunday = true;
          break;
        default:
          toastr.error('Internal Error.');
        }
        i++;
        t.add(1, 'day');
      } while (tString !== dateEndString && i < 1000);
    }
  };

  // weekDayCheck returns true when everything is okay, false otherwise
  $scope.weekDayCheck = function(id) {
    var rule = $scope.reservationRulesById[id];
    var weekDaySelected = rule.Monday || rule.Tuesday || rule.Wednesday ||
          rule.Thursday || rule.Friday || rule.Saturday || rule.Sunday;

    if (weekDaySelected) {
      var okay = !!(rule.DateStart && rule.DateEnd);
      if (!okay) {
        toastr.error('Set Date Start and Date End before setting week days!');
      }
      return okay;
    } else {
      return true;
    }
  };


}]); // app.controller

})(); // closure
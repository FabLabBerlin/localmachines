(function(){

'use strict';
var app = angular.module('fabsmith.admin.activation', ['ngRoute', 'ngCookies', 'fabsmith.admin.randomtoken']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/activations/:activationId', {
    templateUrl: 'ng-modules/activation/activation.html',
    controller: 'ActivationCtrl'
  });
}]); // app.config

app.controller('ActivationCtrl',
 ['$scope', '$routeParams', '$http', '$location', '$cookies', 'randomToken', 'api',
 function($scope, $routeParams, $http, $location, $cookies, randomToken, api) {

  $scope.activation = {
    Id: $routeParams.activationId
  };

  $scope.machines = [];
  $scope.machinesById = {};
  $scope.users = [];
  $scope.usersById = {};

  function loadMachines() {
    api.loadMachines(function(resp) {
      $scope.machines = _.sortBy(resp.machines, function(machine) {
        return machine.Name;
      });
      _.each($scope.machines, function(machine) {
        $scope.machinesById[machine.Id] = machine;
      });
      loadActivation($scope.activation.Id);
    });
  }

  function loadActivation(id) {
    $http({
      method: 'GET',
      url: '/api/activations/' + id,
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(r) {
      $scope.activation = r;
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
      toastr.error('Failed to get activations');
    });
  }

  $scope.refreshUserPicker = function() {
    $('#activation-user').selectpicker('refresh');
  };

  $scope.save = function() {
    $http({
      method: 'PUT',
      url: '/api/activations/' + $scope.activation.Id,
      headers: {'Content-Type': 'application/json' },
      data: $scope.activation,
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
      toastr.success('Activation updated');
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

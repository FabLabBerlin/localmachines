(function(){

'use strict';

var app = angular.module('fabsmith.admin.tutoring', ['ngRoute', 'ngCookies', 'fabsmith.admin.api']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/tutoring', {
    templateUrl: 'ng-modules/tutoring/tutoring.html',
    controller: 'TutoringCtrl'
  });
}]); // app.config

app.controller('TutoringCtrl', ['$scope', '$http', '$location', 'api',
  function($scope, $http, $location, api) {

  $scope.machines = [];
  $scope.tutors = [];

  // Load global settings for the VAT and currency
  $scope.loadSettings = function() {
    $http({
      method: 'GET',
      url: '/api/settings',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(settings) {
      $scope.settings = {
        Currency: {},
        VAT: {}
      };
      _.each(settings, function(setting) {
        $scope.settings[setting.Name] = setting;
      });
    })
    .error(function() {
      toastr.error('Failed to get global config');
    });
  };

  $scope.showTutorSkills = function() {
    if ($scope.machines.length && $scope.tutors.length) {
      // Translate MachineSkills machine ID's to machine names for each tutor
      _.each($scope.tutors, function(tutor) {
        var machineSkills = [];
        if (tutor.MachineSkills) {
          var tmp = tutor.MachineSkills;
          tmp = tmp.slice(1, tmp.length - 1);
          _.each(tmp.split(','), function(idString) {
            var id = parseInt(idString);
            var machine = $scope.machinesById[id];
            if (machine) {
              machineSkills.push(machine);
            }
          });
        }
        tutor.MachineSkills = machineSkills;

      });
    }
  };

  $scope.getMachineNameById = function(machineId) {
    var machineName = '';
    _.each($scope.machines, function(machine) {
      if (parseInt(machine.Id) === parseInt(machineId)) {
        machineName = machine.Name;
      }
    });
    return machineName;
  };

  $scope.loadPurchases = function() {
    $http({
      method: 'GET',
      url: '/api/purchases',
      params: {
        ac: new Date().getTime(),
        type: 'tutor'
      }
    })
    .success(function(response) {
      $scope.purchases = _.sortBy(response.Data, function(purchase) {
        return purchase.Name;
      });
      $scope.purchases = _.map($scope.purchases, function(p) {
        var tutor = $scope.tutorsById[p.ProductId];
        if (tutor) {
          p.Product = tutor;
        }
        p.User = $scope.usersById[p.UserId];
        p.TimeStartLocal = moment(p.TimeStart).tz('Europe/Berlin').format('YYYY-MM-DD HH:mm');
        p.TimeEndLocal = moment(p.TimeEnd).tz('Europe/Berlin').format('YYYY-MM-DD HH:mm');
        return p;
      });
      $scope.showTutorSkills();
    })
    .error(function() {
      toastr.error('Failed to load tutoring purchases');
    });
  };

  $scope.addTutor = function() {
    $http({
      method: 'POST',
      url: '/api/products',
      params: {
        name: name,
        ac: new Date().getTime(),
        type: 'tutor'
      }
    })
    .success(function(data) {
      $scope.editTutor(data.Product.Id);
    })
    .error(function() {
      toastr.error('Failed to create product');
    });
  };

  $scope.editTutor = function(id) {
    $location.path('/tutoring/tutors/' + id);
  };

  $scope.addPurchase = function() {
    $http({
      method: 'POST',
      url: '/api/purchases',
      params: {
        ac: new Date().getTime(),
        type: 'tutor'
      }
    })
    .success(function(tutoringPurchase) {
      $location.path('/tutoring/purchases/' + tutoringPurchase.Id);
    })
    .error(function() {
      toastr.error('Failed to create tutoring purchase');
    });
  };

  $scope.editPurchase = function(id) {
    $location.path('/tutoring/purchases/' + id);
  };

  $scope.loadSettings();
  api.loadMachines(function(machineData) {
    $scope.machines = machineData.machines;
    $scope.machinesById = machineData.machinesById;
    api.loadTutors(function(tutorData) {
      $scope.tutors = tutorData.tutors;
      $scope.tutorsById = tutorData.tutorsById;
      api.loadUsers(function(userData) {
        $scope.users = userData.users;
        $scope.usersById = userData.usersById;
        $scope.loadPurchases();
      });
    });
  });

}]); // app.controller

})(); // closure
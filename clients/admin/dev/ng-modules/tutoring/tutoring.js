(function(){

'use strict';

var app = angular.module('fabsmith.admin.tutoring', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/tutoring', {
    templateUrl: 'ng-modules/tutoring/tutoring.html',
    controller: 'TutoringCtrl'
  });
}]); // app.config

app.controller('TutoringCtrl', ['$scope', '$http', '$location', 
  function($scope, $http, $location) {

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
      console.log(settings);
      _.each(settings, function(setting) {
        $scope.settings[setting.Name] = setting;
      });
    })
    .error(function() {
      toastr.error('Failed to get global config');
    });
  };

  $scope.getAllMachines = function() {
    $http({
      method: 'GET',
      url: '/api/machines',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(machines) {
      $scope.machines = machines;
      $scope.showTutorSkills();
      setTimeout(function() {
        $('.selectpicker').selectpicker('refresh');
      }, 100);
    })
    .error(function(data, status) {
      toastr.error('Failed to get all machines');
    });
  };

  $scope.loadTutors = function() {
    $http({
      method: 'GET',
      url: '/api/tutoring/tutors',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(tutorList) {
      $scope.tutors = tutorList.Data;
      $scope.showTutorSkills();
    })
    .error(function() {
      toastr.error('Failed to load tutor list');
    });
  };

  $scope.showTutorSkills = function() {
    if ($scope.machines.length && $scope.tutors.length) {
      // Translate MachineSkills machine ID's to machine names for each tutor
      _.each($scope.tutors, function(tutor) {
        if (tutor.MachineSkills !== '') {
          var machineSkills = JSON.parse(tutor.MachineSkills);
          tutor.MachineSkills = machineSkills;

          tutor.Skills = '';
          _.each(tutor.MachineSkills, function(machineId, key) {
            tutor.Skills += $scope.getMachineNameById(machineId);
            if (key < tutor.MachineSkills.length - 1) {
              tutor.Skills += ', ';
            }
          });
        }

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
      url: '/api/tutoring/purchases',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(purchaseList) {
      $scope.purchases = purchaseList.Data;
      _.each($scope.purchases, function(purchase) {
        purchase.TutorName = 'Tina Atari';
        purchase.UserName = 'Milda Sane';

        purchase.Created = moment(purchase.StartTime).format('D MMM YY');
        purchase.TimeStart = moment(purchase.StartTime).format('D MMM YY HH:mm');
        purchase.TimeEnd = moment(purchase.EndTime).format('D MMM YY HH:mm');

        purchase.ReservedTimeTotalHours = purchase.TotalTime.toFixed(0);
        purchase.ReservedTimeTotalMinutes = 12;

        purchase.TimerTimeTotalHours = 0;
        purchase.TimerTimeTotalMinutes = 0;

        purchase.TimeTotal = purchase.TotalTime.toFixed(2);
      });
      console.log(purchaseList);
    })
    .error(function() {
      toastr.error('Failed to load purchase list');
    });
  };

  $scope.addTutor = function() {
    $location.path('/tutoring/tutor');
  };

  $scope.editTutor = function() {
    $location.path('/tutoring/tutor');
  };

  $scope.addPurchase = function() {
    $location.path('/tutoring/purchase');
  };

  $scope.editPurchase = function() {
    $location.path('/tutoring/purchase');
  };

  $scope.loadSettings();
  $scope.getAllMachines();
  $scope.loadTutors();
  $scope.loadPurchases();

}]); // app.controller

})(); // closure
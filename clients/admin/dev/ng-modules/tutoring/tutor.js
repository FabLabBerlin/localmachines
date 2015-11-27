(function(){

'use strict';

var app = angular.module('fabsmith.admin.tutoring.tutor', ['ngRoute', 'ngCookies', 'fabsmith.admin.api']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/tutoring/tutor', {
    templateUrl: 'ng-modules/tutoring/tutor.html',
    controller: 'TutorCtrl'
  });
}]); // app.config

app.controller('TutorCtrl', ['$scope', '$http', '$location', 'api',
  function($scope, $http, $location, api) {

  $scope.users = [];
  $scope.machines = [];
  $scope.tutor = {
    PriceUnit: 'hour',
    MachineSkills: []
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
      setTimeout(function() {
        $('.selectpicker').selectpicker('refresh');
      }, 100);
    })
    .error(function(data, status) {
      toastr.error('Failed to get all machines');
    });
  };

  api.loadUsers(function(usersData) {
    $scope.users = usersData.users;
    setTimeout(function() {
      $('.selectpicker').selectpicker('refresh');
    }, 100);
  });
  $scope.getAllMachines();

  $scope.machineSkillAdded = function(machineSkill) {
    var skillFound = false;
    for (var i=0; i<$scope.tutor.MachineSkills.length; i++) {
      if (parseInt($scope.tutor.MachineSkills[i].Id) ===
        parseInt(machineSkill.Id)) {
        skillFound = true;
        break;
      }
    }
    return skillFound;
  };

  $scope.addMachineSkill = function() {
    $('#skill-picker').selectpicker('refresh');
    
    // Get machine by id
    for (var i=0; i<$scope.machines.length; i++) {
      if (parseInt($scope.SelectedMachineId) === 
        parseInt($scope.machines[i].Id)) {
        
        // Check for existing skill
        var skill = {
          Id: $scope.machines[i].Id, 
          Name: $scope.machines[i].Name
        };

        // Add only if skill is not there yet
        if (!$scope.machineSkillAdded(skill)) {
          $scope.tutor.MachineSkills.push(skill);
        } else {
          toastr.warning('This skill is already added');
        }
      }
    }
  };

  $scope.removeMachineSkill = function(skillId) {
    for (var i=0; i<$scope.tutor.MachineSkills.length; i++) {
      if (parseInt(skillId) === parseInt($scope.tutor.MachineSkills[i].Id)) {
        $scope.tutor.MachineSkills.splice(i, 1);
      }
    }
  };

  $scope.cancel = function() {
    $location.path('/tutoring');
  };

  $scope.save = function() {

    if (!$scope.tutor.UserId) {
      toastr.error('Select tutor user');
      return;
    }

    if ($scope.tutor.Price === '' || !$scope.tutor.Price) {
      toastr.error('Enter tutor price');
      return;
    }

    if (isNaN($scope.tutor.Price)) {
      toastr.error('Tutor price should be a number');
      return;
    }

    $scope.updateTutor();
  };

  $scope.updateTutor = function() {
    $http({
      method: 'PUT',
      url: '/api/tutoring/tutor',        
      headers: {'Content-Type': 'application/json' },
      data: $scope.tutor,
      transformRequest: function(tutor) {

        var transformedTutor = {
          UserId: parseInt(tutor.UserId),
          Price: parseFloat(tutor.Price),
          PriceUnit: tutor.PriceUnit,
          MachineSkills: ''
        };

        var tutorSkills = '[';
        for (var i=0; i<tutor.MachineSkills.length; i++) {
          tutorSkills += tutor.MachineSkills[i].Id;
          if (i < tutor.MachineSkills.length - 1) {
            tutorSkills += ',';
          }
        }
        tutorSkills += ']';

        transformedTutor.MachineSkills = tutorSkills;
        console.log(transformedTutor);

        var container = {
          Product: transformedTutor
        };

        return JSON.stringify(container);
      },
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(updatedTutor) {
      console.log(updatedTutor);

      // Udpdate the id of the tutor if created
      $scope.tutor.Id = updatedTutor.Product.Id;
      $scope.tutor.Name = updatedTutor.Product.Name;
      
      $location.path('/tutoring');
    })
    .error(function(data) {
      console.log(data);
      toastr.error('Failed to update tutor');
    });
  };

}]); // app.controller

})(); // closure
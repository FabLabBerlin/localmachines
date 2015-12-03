(function(){

'use strict';

var app = angular.module('fabsmith.admin.tutoring.tutor', ['ngRoute', 'ngCookies', 'fabsmith.admin.api']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/tutoring/tutors/:id', {
    templateUrl: 'ng-modules/tutoring/tutor.html',
    controller: 'TutorCtrl'
  });
}]); // app.config

app.controller('TutorCtrl', ['$scope', '$routeParams', '$http', '$location', 'api',
  function($scope, $routeParams, $http, $location, api) {

  $scope.users = [];
  $scope.machines = [];
  $scope.machinesById = {};
  $scope.tutor = {
    Product: {
      Id: $routeParams.id,
      PriceUnit: 'hour',
      MachineSkills: []
    }
  };

  function loadTutor() {
    $http({
      method: 'GET',
      url: '/api/products/' + $scope.tutor.Product.Id,
      params: {
        ac: new Date().getTime(),
        type: 'tutor'
      }
    })
    .success(function(tutor) {
      $scope.tutor = tutor;
      var machineSkills = [];
      if ($scope.tutor.Product.MachineSkills) {
        var tmp = $scope.tutor.Product.MachineSkills;
        tmp = tmp.slice(1, tmp.length - 1);
        _.each(tmp.split(','), function(idString) {
          var id = parseInt(idString);
          var machine = $scope.machinesById[id];
          if (machine) {
            machineSkills.push(machine);
          }
        });
      }
      $scope.tutor.Product.MachineSkills = machineSkills;
      console.log('$scope.tutor = ', $scope.tutor);
    })
    .error(function(data, status) {
      toastr.error('Failed to load tutor data');
    });
  }

  $scope.refreshUserPicker = function() {
    $('#user-picker').selectpicker('refresh');
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
      $scope.machinesById = {};
      _.each(machines, function(machine) {
        $scope.machinesById[machine.Id] = machine;
      });

      setTimeout(function() {
        $('.selectpicker').selectpicker('refresh');
      }, 100);
      loadTutor();
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

  $scope.machineSkillAdded = function(machineSkill) {
    var skillFound = false;
    for (var i=0; i<$scope.tutor.Product.MachineSkills.length; i++) {
      if (parseInt($scope.tutor.Product.MachineSkills[i].Id) ===
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
          $scope.tutor.Product.MachineSkills.push(skill);
        } else {
          toastr.warning('This skill is already added');
        }
      }
    }
  };

  $scope.removeMachineSkill = function(skillId) {
    for (var i=0; i<$scope.tutor.Product.MachineSkills.length; i++) {
      if (parseInt(skillId) === parseInt($scope.tutor.Product.MachineSkills[i].Id)) {
        $scope.tutor.Product.MachineSkills.splice(i, 1);
      }
    }
  };

  $scope.cancel = function() {
    $location.path('/tutoring');
  };

  $scope.save = function() {

    if (!$scope.tutor.Product.UserId) {
      toastr.error('Select tutor user');
      return;
    }

    if ($scope.tutor.Product.Price === '' || !$scope.tutor.Product.Price) {
      toastr.error('Enter tutor price');
      return;
    }

    if (isNaN($scope.tutor.Product.Price)) {
      toastr.error('Tutor price should be a number');
      return;
    }

    $scope.updateTutor();
  };

  $scope.updateTutor = function() {
    $http({
      method: 'PUT',
      url: '/api/products/' + $scope.tutor.Product.Id + '?type=tutor',
      headers: {'Content-Type': 'application/json' },
      data: $scope.tutor,
      transformRequest: function(tutor) {
        var transformed = {
          Product: _.extend({}, tutor.Product)
        };
        transformed.Product.Id = parseInt(tutor.Product.Id);
        transformed.Product.UserId = parseInt(tutor.Product.UserId);
        transformed.Product.Price = parseFloat(tutor.Product.Price);

        var tutorSkills = '[';
        tutorSkills += _.pluck(tutor.Product.MachineSkills, 'Id').join(',');
        tutorSkills += ']';

        transformed.Product.MachineSkills = tutorSkills;
        console.log(transformed);

        return JSON.stringify(transformed);
      },
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(updatedTutor) {
      $location.path('/tutoring');
    })
    .error(function(data) {
      console.log(data);
      toastr.error('Failed to update tutor');
    });
  };

  api.loadMachines(function(data) {
    $scope.machines = data.machines;
    $scope.machinesById = data.machinesById;

    setTimeout(function() {
      $('.selectpicker').selectpicker('refresh');
    }, 100);
    loadTutor();
  });

}]); // app.controller

})(); // closure
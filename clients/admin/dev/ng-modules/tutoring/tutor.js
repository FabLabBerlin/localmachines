(function(){

'use strict';

var app = angular.module('fabsmith.admin.tutoring.tutor', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/tutoring/tutor', {
    templateUrl: 'ng-modules/tutoring/tutor.html',
    controller: 'TutorCtrl'
  });
}]); // app.config

app.controller('TutorCtrl', ['$scope', '$http', '$location', 
  function($scope, $http, $location) {

  $scope.users = [];
  $scope.machines = [];
  $scope.tutor = {
    PriceUnit: 'hour',
    Skills: []
  };
  
  $scope.getAllUsers = function() {
    $http({
      method: 'GET',
      url: '/api/users',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(users) {
      $scope.users = users;
      setTimeout(function() {
        $('.selectpicker').selectpicker('refresh');
      }, 100);
    })
    .error(function(data, status) {
      toastr.error('Failed to get all users');
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
      setTimeout(function() {
        $('.selectpicker').selectpicker('refresh');
      }, 100);
    })
    .error(function(data, status) {
      toastr.error('Failed to get all machines');
    });
  };

  $scope.getAllUsers();
  $scope.getAllMachines();

  $scope.skillAdded = function(skill) {
    var skillFound = false;
    for (var i=0; i<$scope.tutor.Skills.length; i++) {
      if (parseInt($scope.tutor.Skills[i].Id) ===
        parseInt(skill.Id)) {
        skillFound = true;
        break;
      }
    }
    return skillFound;
  };

  $scope.addSkill = function() {
    $('#skill-picker').selectpicker('refresh');
    
    // Get machine by id
    for (var i=0; i<$scope.machines.length; i++) {
      if (parseInt($scope.tutor.SelectedMachineId) === 
        parseInt($scope.machines[i].Id)) {
        
        // Check for existing skill
        var skill = {
          Id: $scope.machines[i].Id, 
          Name: $scope.machines[i].Name
        };

        // Add only if skill is not there yet
        if (!$scope.skillAdded(skill)) {
          console.log('Adding skill: ' + $scope.tutor.SelectedMachineId);
          $scope.tutor.Skills.push(skill);
        }
      }
    }
  };

  $scope.cancel = function() {
    $location.path('/tutoring');
  };

  $scope.save = function() {

    if ($scope.tutor.Name === '' || !$scope.tutor.Name) {
      alert('Enter tutor name');
      return;
    }

    if ($scope.tutor.Price === '' || !$scope.tutor.Price) {
      alert('Enter tutor price');
      return;
    }

    console.log($scope.tutor);

    var tutor = $scope.tutor;
    tutor.Price = parseInt(tutor.Price);

    /*
    $http({
      method: 'PUT',
      url: '/api/tutoring/tutor',
      data: tutor,
      headers: { 
        'Content-Type': 'application/json' 
      },
      params: { 
        ac: new Date().getTime() 
      },
      transformRequest: function(data) {
        return JSON.stringify(data);
      }
    })
    .success(function(data) {
      toastr.success('Tutor updated');
      $location.path('/tutoring');
    })
    .error(function() {
      toastr.error('Failed to save tutor');
    });
    */

    $scope.updateTutor = function() {
      if (!$scope.tutor.Product.Id) {
        // The backend should create a new product if the Id is 0
        $scope.tutor.Product.Id = 0;
      }

      $http({
        method: 'PUT',
        url: '/api/products/' + $scope.tutor.Product.Id + '?type=tutor',
        headers: {'Content-Type': 'application/json' },
        data: $scope.space,
        transformRequest: function(data) {
          var transformed = {
            Product: _.extend({}, data.Product)
          };
          transformed.Product.Id = parseInt(data.Product.Id);
          transformed.Product.Price = parseFloat(data.Product.Price);
          return JSON.stringify(transformed);
        },
        params: {
          ac: new Date().getTime()
        }
      })
      .success(function(data) {
        toastr.success('Update successful');
      })
      .error(function(data) {
        console.log(data);
        toastr.error('Failed to update');
      });
    };

  };

}]); // app.controller

})(); // closure
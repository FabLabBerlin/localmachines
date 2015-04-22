(function(){

'use strict';

// Constants for machine states
var MACHINE_STATUS_AVAILABLE = 'free';
var MACHINE_STATUS_OCCUPIED = 'occupied';
var MACHINE_STATUS_USED = 'used'; // Used by the current user
var MACHINE_STATUS_UNAVAILABLE = 'unavailable';

// Our local app variable for the module

var app = angular.module('fabsmith.machines', 
  ['ngRoute', 'fabsmithFilters']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/machines', {
    templateUrl: 'ng-modules/machines/machines.html',
    controller: 'MachinesCtrl'
  });
}]);

app.controller('MachinesCtrl', 
 ['$scope', '$http', '$location', '$route', '$cookieStore', 
 function($scope, $http, $location, $route, $cookieStore) {

  // Attempt to get logged user machines
  $scope.loadMachines = function() {
    $http({
      method: 'GET',
      url: '/api/users/' + $cookieStore.get('Id') + '/machines'
    })
    .success(function(machines){
      console.log(machines);
      $scope.onMachinesLoaded(machines);
    })
    .error(function(data, status) {
      if (status === 401) {
        alert('Not authorized');
      } else {
        alert('Error loading machines');
      }
    });
  };

  $scope.userFullName = $cookieStore.get('FirstName') + 
    ' ' + $cookieStore.get('LastName');
  $scope.machines = [];
  $scope.loadMachines();

  $scope.getActiveActivations = function(machines) {
    $http({
      method: 'GET',
      url: '/api/activations/active'
    })
    .success(function(activations){
      console.log(activations);
      $scope.onActivationsLoaded(activations, machines);
    })
    .error(function(data, status){
      alert('Failed to load active activations');
    });
  };

  $scope.onMachinesLoaded = function(machines){
    if (machines.length <= 0) {
      alert('There are no machines available for you');
    } else if (machines.length > 0) {
      $scope.getActiveActivations(machines);
    } else {
      alert('Error loading machines');
    }
  };

  $scope.onActivationsLoaded = function(activations, machines){

    // Got activations
    // Add status vars to machines
    for (var mchIter = 0; mchIter < machines.length; mchIter++) {
      
      // TODO: figure out simpler way for indicating machine status
      machines[mchIter].available = false;
      machines[mchIter].used = false;
      machines[mchIter].occupied = false;
      machines[mchIter].unavailable = true;

      if (machines[mchIter].Available) {
        machines[mchIter].available = true;
        machines[mchIter].unavailable = false;
      } else {
        
        // If machine is not available it means that
        // it is either occupied by someone else, 
        // unavailable or used by the user logged
        for (var actIter = 0; actIter < activations.length; actIter++) {
          if (activations[actIter].MachineId === machines[mchIter].Id) {
            
            var activationStartTime;
            var timeNow;
            var activationElapsedTime;

            if ($cookieStore.get('Id') === activations[actIter].UserId) {

              // Machine is being used by logged user
              machines[mchIter].used = true;
              machines[mchIter].unavailable = false;

              // Assign other activation related data to the machine object
              machines[mchIter].OccupiedByUserId = activations[actIter].UserId;
              machines[mchIter].ActivationId = activations[actIter].Id;

              // What we also need is to start the counter interval
              // Start timer for elapsed time
              machines[mchIter].ActivationSecondsElapsed = 
                  $scope.getActivationElapsedSeconds(activations[actIter]);
              machines[mchIter].activationInterval = 
                setInterval($scope.updateElapsedTime, 1000, mchIter);

            } else {

              // Machine is being used by someone else
              machines[mchIter].occupied = true;
              machines[mchIter].unavailable = false;
              machines[mchIter].OccupiedByUserId = activations[actIter].UserId;
              machines[mchIter].ActivationId = activations[actIter].Id;
              machines[mchIter].UserAdmin = 
                ($cookieStore.get('UserRole') === 'admin');
              
              // But, if logged as admin, we can also set the activation ID
              // and elapsed time
              if ( machines[mchIter].UserAdmin ) {
                machines[mchIter].ActivationSecondsElapsed = 
                  $scope.getActivationElapsedSeconds(activations[actIter]);
                machines[mchIter].activationInterval = 
                  setInterval($scope.updateElapsedTime, 1000, mchIter);
              }

              // Get user name
              $scope.getOccupierName(machines[mchIter], activations[actIter].UserId);
            }
          }
        } // for activations
      } // if machine available else
    } // for machines

    $scope.machines = machines;
  };

  $scope.getActivationElapsedSeconds = function(activation){
    var activationStartTime = Date.parse(activation.TimeStart);
    var timeNow = Date.now();
    var activationElapsedTime = timeNow - activationStartTime;
    activationElapsedTime = Math.round(activationElapsedTime / 1000);
    return activationElapsedTime;
  };

  $scope.getOccupierName = function(machine, userId) {
    $http({
      method: 'GET',
      url: '/api/users/' + userId + '/name'
    })
    .success(function(data){
      machine.occupier = data.FirstName + ' ' + data.LastName;
    })
    .error(function(){
      toastr.error('Failed to get occupier name');
    });
  };

  $scope.showGlobalLoader = function() {
    $('#loader-global').removeClass('hidden');
  };

  $scope.hideGlobalLoader = function() {
    $('#loader-global').addClass('hidden');
  };

  $scope.updateElapsedTime = function(machineIter) {
    $scope.machines[machineIter].ActivationSecondsElapsed++;
    $scope.$apply();
  };

  // Activate a machine by the currenty logged in user
  $scope.activate = function(machineId) {
    $scope.showGlobalLoader();
    $http({
      method: 'POST',
      url: '/api/activations',
      params: {
        mid: machineId
      }
    })
    .success(function(data) {
      console.log(data);
      $scope.hideGlobalLoader();

      // Find machine by ID
      for (var machineIter = 0; machineIter < $scope.machines.length; machineIter++) {

        // Machine found condition
        if ($scope.machines[machineIter].Id === machineId) {
          
          // Refresh data of the new activation
          $scope.machines[machineIter].ActivationSecondsElapsed = 0;
          $scope.machines[machineIter].OccupiedByUserId = parseInt( $cookieStore.get('Id') );
          $scope.machines[machineIter].ActivationId = data.ActivationId;
          $scope.machines[machineIter].available = false;
          $scope.machines[machineIter].used = true;
  
          // Start timer for elapsed time
          $scope.machines[machineIter].activationInterval = 
            setInterval($scope.updateElapsedTime, 1000, machineIter);
  
          // Exit the machine finding for loop
          break;
        }
      }
    })
    .error(function() {
      $scope.hideGlobalLoader();
      toastr.error('Could not activate machine');
    });

  };

  $scope.deactivatePrompt = function(machine) {
    vex.dialog.buttons.YES.text = 'Yes';
    vex.dialog.buttons.NO.text = 'No';
    vex.dialog.confirm({
      message: 'Are you sure?',
      callback: $scope.deactivatePromptCallback.bind(this, machine)
    });
  };

  $scope.deactivatePromptCallback = function(machine, value) {
    if (value) {    
      $scope.deactivate(machine);
    }
  };

  $scope.deactivate = function(machine) {

    // Stop activation timer interval
    clearInterval(machine.activationInterval);

    // Show loader
    $scope.showGlobalLoader();

    $http({
      method: 'PUT', 
      url: '/api/activations/' + machine.ActivationId
    })
    .success(function(data) {
      $scope.hideGlobalLoader();
      machine.used = false;
      machine.occupied = false;
      machine.available = true;
    })
    .error(function() {
      $scope.hideGlobalLoader();
      toastr.error('Failed to deactivate');
    });
  };

}]);

app.directive('fsMachineItem', function() {
  return {
    templateUrl: 'ng-modules/machines/machine-item.html',
    restrict: 'E',
    controller: ['$scope', '$element', function($scope, $element) {
      
      $scope.infoVisible = false;

      if ($scope.infoVisible) {
        $($element).find('.machine-info-content').first().show();
      } else {
        $($element).find('.machine-info-content').first().hide();
      }

      $scope.toggleInfo = function() {
        $scope.infoVisible = !$scope.infoVisible;
        if ($scope.infoVisible) {
          $($element).find('.machine-info-content').first().slideDown();
        } else {
          $($element).find('.machine-info-content').first().slideUp();
        }
      };

    }]
  };
});

app.directive('fsMachineBodyAvailable', function() {
  return {
    templateUrl: 'ng-modules/machines/machine-body-available.html',
    restrict: 'E'
  };
});

app.directive('fsMachineBodyUsed', function() {
  return {
    templateUrl: 'ng-modules/machines/machine-body-used.html',
    restrict: 'E',
    controller: ['$scope', function($scope){
      
      if ($scope.machine.Status === MACHINE_STATUS_USED) {
        $scope.machine.activationInterval = setInterval(function() {
          $scope.machine.ActivationSecondsElapsed++;
          $scope.$apply();
        }, 1000);
      }

    }]
  };
});

app.directive('fsMachineBodyOccupied', function() {
  return {
    templateUrl: 'ng-modules/machines/machine-body-occupied.html',
    restrict: 'E',
    controller: ['$scope', '$cookieStore', function($scope, $cookieStore){

      // As we are using this scope for more than one directive
      if ($scope.machine.occupied) {
        var user = {};
        user.Admin = $cookieStore.get('UserRole') === 'admin';
        $scope.user = user;

        // Activate occupied machine timer if user is admin or staff
        if (user.Admin) {
          $scope.machine.activationInterval = setInterval(function() {
            $scope.machine.ActivationSecondsElapsed++;
            $scope.$apply();
          }, 1000);
        }

      }

    }]
  };
});

app.directive('fsMachineBodyUnavailable', function() {
  return {
    templateUrl: 'ng-modules/machines/machine-body-unavailable.html',
    restrict: 'E'
  };
});

})(); // closure
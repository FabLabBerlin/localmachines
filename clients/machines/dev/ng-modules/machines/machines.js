(function(){

'use strict';

// Constants for machine states
var MACHINE_STATUS_AVAILABLE = 'free';
var MACHINE_STATUS_OCCUPIED = 'occupied';
var MACHINE_STATUS_USED = 'used'; // Used by the current user
var MACHINE_STATUS_UNAVAILABLE = 'unavailable';
var LOGOUT_TIMER_DELAY = 30;

// Our local app variable for the module

var app = angular.module('fabsmith.machines', 
  ['ngRoute', 'timer', 'fabsmithFilters', 'ui.bootstrap.modal']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/machines', {
    templateUrl: 'ng-modules/machines/machines.html',
    controller: 'MachinesCtrl'
  });
}]);

app.controller('MachinesCtrl', 
 ['$scope', '$http', '$location', '$route', '$cookieStore', '$modal', 
 function($scope, $http, $location, $route, $cookieStore, $modal) {

  // Attempt to get logged user machines
  $scope.loadMachines = function() {
    $http({
      method: 'GET',
      url: '/api/users/' + $cookieStore.get('Id') + '/machines',
      params: {
        anticache: new Date().getTime()
      }
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

  $scope.userFullName = $cookieStore.get('FirstName') + ' ' + $cookieStore.get('LastName');
  $scope.machines = [];
  $scope.loadMachines();

  $scope.getActiveActivations = function(machines) {
    $http({
      method: 'GET',
      url: '/api/activations/active',
      params: {
        anticache: new Date().getTime()
      }
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

  // Configure timer
  $scope.resetTimer = function() {
    $scope.$broadcast('timer-clear');
    $scope.$broadcast('timer-set-countdown', LOGOUT_TIMER_DELAY);
    $scope.$broadcast('timer-start');
  };
  $scope.resetTimer();

  $scope.getOccupierName = function(machine, userId) {
    $http({
      method: 'GET',
      url: '/api/users/' + userId + '/name',
      params: {
        anticache: new Date().getTime()
      }
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
    $scope.resetTimer();
    $scope.showGlobalLoader();
    $http({
      method: 'POST',
      url: '/api/activations',
      params: {
        mid: machineId,
        anticache: new Date().getTime()
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
    
    $scope.resetTimer();
    
    var modal = $scope.openDeactivateModal();
    modal.result.then(function() {
      $scope.resetTimer();
      $scope.deactivate(machine);
    }, function () {
      $scope.resetTimer();
    });
  };

  $scope.deactivate = function(machine) {

    // Stop activation timer interval
    clearInterval(machine.activationInterval);

    // Show loader
    $scope.showGlobalLoader();

    $http({
      method: 'PUT', 
      url: '/api/activations/' + machine.ActivationId, 
      params: {
        anticache: new Date().getTime()
      }
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

  // TODO: use the activation api calls to handle the following
  $scope.isAvailable = function(machine) {
    return machine.Status === MACHINE_STATUS_AVAILABLE;
  };

  $scope.isOccupied = function(machine) {
    return machine.Status === MACHINE_STATUS_OCCUPIED;
  };

  $scope.isUsed = function(machine) {
    return machine.Status === MACHINE_STATUS_USED;
  };

  $scope.isUnavailable = function(machine) {
    return machine.Status === MACHINE_STATUS_UNAVAILABLE;
  };

  $scope.setAllStates = function(machine, trueOrFalse) {
    machine.available = trueOrFalse;
    machine.used = trueOrFalse;
    machine.occupied = trueOrFalse;
    machine.unavailable = trueOrFalse;
  };

  // TODO: Remove the angular-ui-bootstrap dependency ans substitute with 
  //       plain Bootstrap HTML. Current solution causes path problems.
  $scope.openDeactivateModal = function() {
    var modalInstance = $modal.open({
      backdrop: false,
      templateUrl: 'ng-modules/machines/deactivate-modal.html?v1',
      windowTemplateUrl: '/views/bower_components/angular-ui-bootstrap/template/modal/window.html',
      controller: 'DeactivateModalCtrl'
      });

      return modalInstance;

  }; // showModal

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

        $scope.resetTimer();

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
        user.Admin = $cookieStore.get('Admin');
        user.Staff = $cookieStore.get('Staff');
        user.Member = $cookieStore.get('Member');
        $scope.user = user;

        // Activate occupied machine timer if user is admin or staff
        if (user.Admin || user.Staff) {
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

app.controller('DeactivateModalCtrl', function ($scope, $modalInstance) {

  $scope.proceed = function () {
    $modalInstance.close('eh');
  };

  $scope.cancel = function () {
    $modalInstance.dismiss('cancel');
  };

});

})(); // closure
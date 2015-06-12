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

  $scope.scrollTop = false;
  $scope.scrollBottom = false;

  if (window.libnfc) {

    $scope.onNfcError = function(error) {
      window.libnfc.cardRead.disconnect($scope.nfcLogout);
      window.libnfc.cardReaderError.disconnect($scope.onNfcError);
      toastr.error(error);
      $scope.nfcTimeout = setTimeout($scope.activateNfcLogout, 2000);
    };

    $scope.onNfc = function(uid) {
      $scope.smartLogout();
    };

    $scope.activateNfcLogout = function() {
      window.libnfc.cardRead.connect($scope.onNfc);
      window.libnfc.cardReaderError.connect($scope.onNfcError);
      window.libnfc.asyncScan();
      toastr.info('You can log out by using your NFC card');
    };

    // NFC logout functionality is activated N seconds after log in
    $scope.nfcTimeout = setTimeout($scope.activateNfcLogout, 1000);

    setTimeout(function(){
      $scope.scrollTop = true;
      $scope.scrollBottom = true;
      $('.scroll-nav-top').hide();
      $scope.currentScroll = 0;
      $scope.scrollStep = $(window).height() / 2;
      $scope.checkScroll();
    }, 200);

    $scope.checkScroll = function() {
      if($(window).height() < $("html, body").height()){
        if ($(window).height() + $scope.currentScroll >= $('html,body').height()) {
          //$scope.scrollBottom = false;
          $('.scroll-nav-bottom').slideUp();
        } else {
          $('.scroll-nav-bottom').slideDown();
        }

        if ($scope.currentScroll <= 0) {
          $('.scroll-nav-top').slideUp();
        } else {
          $('.scroll-nav-top').slideDown();
        }
      } else {
        $('.scroll-nav-top').hide();
        $('.scroll-nav-bottom').hide();
      }
    };

    $scope.scrollUp = function() {
      $scope.currentScroll -= $scope.scrollStep;
      $scope.checkScroll();
      if ( $scope.currentScroll <= 0 ) {
        $scope.currentScroll = 0;
      }
      $('html,body').animate({
        scrollTop: $scope.currentScroll
      }, 'easeOutExpo');
    };

    $scope.scrollDown = function() {
      $scope.currentScroll += $scope.scrollStep;
      $scope.checkScroll();
      if ($scope.currentScroll + $(window).height() >= $('html,body').height()) {
        $scope.currentScroll = $('html,body').height() - $(window).height();
      }
      $('html,body').animate({
        scrollTop: $scope.currentScroll
      }, 'easeOutExpo');
    };

    $scope.watchingActivity = false;
    $scope.numberOfSecondsBeforeLogOut = 30;
    (function(nbSecLogOut){
      var idleTime = 0;
      $(document).ready(function () {
        if(!$scope.watchingActivity){
          $scope.watchingActivity = true;
          $scope.idleInterval = setInterval(timerIncrement, 1000);

          $(this).mousemove(function (e) {
            idleTime = 0;
          });
          $(this).keypress(function (e) {
            idleTime = 0;
          });
        }
      });

      function timerIncrement() {
        idleTime = idleTime + 1;
        if (idleTime > (nbSecLogOut-1)) {
          toastr.info('Inactivity time out reached. If you want to manage the machines you can log in again.');
          $scope.smartLogout();
        }
      }
    })($scope.numberOfSecondsBeforeLogOut);
  }

  // Reloading machines every 1 second
  $scope.reloadMachinesInterval = setInterval(function(){
    $scope.loadMachines();
  }, 1000);

  // Makes sure that NFC intervals are cleared
  $scope.smartLogout = function() {
    // Disconnecting events from the nfc reader
    if(window.libnfc) {
      window.libnfc.cardRead.disconnect($scope.onNfc);
      window.libnfc.cardReaderError.disconnect($scope.onNfcError);
    }
    // Clearing timers of all machines
    $.map($scope.machines, function(machine, i){
      if(machine.activationInterval !== undefined){
        clearInterval(machine.activationInterval);
      }
    });
    clearInterval($scope.reloadMachinesInterval);
    clearInterval($scope.idleInterval);
    clearTimeout($scope.nfcTimeout);
    $scope.logout();
  };

  // Attempt to get logged user machines
  $scope.loadMachines = function() {
    $http({
      method: 'GET',
      url: '/api/users/' + $cookieStore.get('Id') + '/machines',
      params: { ac: new Date().getTime() }
    })
    .success(function(machines){
      $scope.onMachinesLoaded(machines);
    })
    .error(function(data, status) {
      toastr.error('Failed to load machines');
      $scope.smartLogout();
    });
  };

  $scope.userFullName = $cookieStore.get('FirstName') +
    ' ' + $cookieStore.get('LastName');
  $scope.machines = [];
  $scope.loadMachines();

  $scope.getActiveActivations = function(machines) {
    $http({
      method: 'GET',
      url: '/api/activations/active',
      params: { ac: new Date().getTime() }
    })
    .success(function(activations){
      $scope.onActivationsLoaded(activations, machines);
    })
    .error(function(data, status){
      toastr.error('Loading activations failed');
      $scope.smartLogout();
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
    $scope.machines = _.map(machines, function(machine, i) {
      if (machine.Image) {
        machine.ImageUrl = '/files/' + machine.Image;
      }

      // TODO: figure out simpler way for indicating machine status
      machine.available = false;
      machine.used = false;
      machine.occupied = false;
      machine.unavailable = true;

      if (machine.Available) {
        machine.available = true;
        machine.unavailable = false;
      } else {

        // If machine is not available it means that
        // it is either occupied by someone else,
        // unavailable or used by the user logged
        for (var actIter = 0; actIter < activations.length; actIter++) {
          if (activations[actIter].MachineId === machine.Id) {

            var activationStartTime;
            var timeNow;
            var activationElapsedTime;

            if ($cookieStore.get('Id') === activations[actIter].UserId) {

              // Machine is being used by logged user
              machine.used = true;
              machine.unavailable = false;

              // Assign other activation related data to the machine object
              machine.OccupiedByUserId = activations[actIter].UserId;
              machine.ActivationId = activations[actIter].Id;

              // What we also need is to start the counter interval
              // Start timer for elapsed time
              machine.ActivationSecondsElapsed =
                  $scope.getActivationElapsedSeconds(activations[actIter]);

            } else {

              // Machine is being used by someone else
              machine.occupied = true;
              machine.unavailable = false;
              machine.OccupiedByUserId = activations[actIter].UserId;
              machine.ActivationId = activations[actIter].Id;
              machine.UserAdmin =
                ($cookieStore.get('UserRole') === 'admin');

              // But, if logged as admin, we can also set the activation ID
              // and elapsed time
              if ( machine.UserAdmin ) {
                machine.ActivationSecondsElapsed =
                  $scope.getActivationElapsedSeconds(activations[actIter]);
              }

              // Get user name
              $scope.getOccupierName(machine, activations[actIter].UserId);
            }
          }
        } // for activations
      } // if machine available else
    }); // for machines

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
      url: '/api/users/' + userId + '/name',
      params: { ac: new Date().getTime() }
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
        mid: machineId,
        ac: new Date().getTime()
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
      message: 'Do you really want to stop the activation for <b>' +
        machine.Name + '</b>?',
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
      url: '/api/activations/' + machine.ActivationId,
      params: { ac: new Date().getTime() }
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
      if ($scope.machine.used) {
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

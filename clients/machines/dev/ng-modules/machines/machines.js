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

	// Show logged user name
	$scope.userFullName = $cookieStore.get('FirstName') + ' ' + $cookieStore.get('LastName');

	// Configure timer
	$scope.resetTimer = function() {
		$scope.$broadcast('timer-clear');
		$scope.$broadcast('timer-set-countdown', LOGOUT_TIMER_DELAY);
		$scope.$broadcast('timer-start');
	};
	$scope.resetTimer();

	$scope.showGlobalLoader = function() {
		$('#loader-global').removeClass('hidden');
	};

	$scope.hideGlobalLoader = function() {
		$('#loader-global').addClass('hidden');
	};

	// Initialize the machines array
	$scope.machines = [];

	$scope.getOccupierName = function(machine, userId) {
		$http({
			method: 'GET',
			url: '/api/users/' + userId + '/fullname',
			params: {
				anticache: new Date().getTime()
			}
		})
		.success(function(data){
			machine.occupier = data.FullName;
		})
		.error(function(){
			alert('Failed to get occupier name');
		});
	};

	// Get current user machines
	$http({
		method: 'GET',
		url: '/api/users/' + $cookieStore.get('Id') + '/machines',
		params: {
			anticache: new Date().getTime()
		}
	})
	.success(function(machines) {
		console.log(machines);
		if (machines.length <= 0) {
			alert('There are no machines available for you');
		} else if (machines.length > 0) {

		// Got machines
		// Get activations
		$http({
			method: 'GET',
			url: '/api/activations/active',
			params: {
				anticache: new Date().getTime()
			}
		})
		.success(function(activations){
			
			// Got activations
			// Add status vars to machines
			for (var machinesIter = 0; machinesIter < machines.length; machinesIter++) {
				
				// TODO: figure out simpler way for indicating machine status
				machines[machinesIter].available = false;
				machines[machinesIter].used = false;
				machines[machinesIter].occupied = false;
				machines[machinesIter].unavailable = true;

				if (machines[machinesIter].Available) {
					machines[machinesIter].available = true;
					machines[machinesIter].unavailable = false;
				} else {
					
					// If machine is not available it means that
					// it is either occupied by someone else, unavailable or used by the user logged
					for (var activationsIter = 0; activationsIter < activations.length; activationsIter++) {
						if (activations[activationsIter].MachineId === machines[machinesIter].Id) {
							if ($cookieStore.get('Id') === activations[activationsIter].UserId) {

								// Machine is being used by logged user
								machines[machinesIter].used = true;
								machines[machinesIter].unavailable = false;
							} else {

								// Machine is being used by someone else
								machines[machinesIter].occupied = true;
								machines[machinesIter].unavailable = false;

								// Get user name
								$scope.getOccupierName(machines[machinesIter], activations[activationsIter].UserId);
							}

						}
					} // for activations
				} // if machine available else
			} // for machines

			$scope.machines = machines;
		})
		.error(function(data, status){
			alert('Failed to load active activations');
		});
		} else {
			alert('Error loading machines');
		}
	})
	.error(function(data, status) {
		if (status === 401) {
			alert('Not authorized');
		} else {
			alert('Error loading machines');
		}
	});

	$scope.updateElapsedTime = function(machineIter) {
		$scope.machines[machineIter].ActivationSecondsElapsed++;
		$scope.$apply();
	};

	// Activate a machine by the currenty logged in user
	$scope.activate = function(machineId) {
		
		$scope.resetTimer();

		// Show loader
		$scope.showGlobalLoader();

		$http({
			method: 'POST',
			url: '/api/activations',
			params: {
				machine_id: machineId,
				anticache: new Date().getTime()
			}
		})
		.success(function(data) {

			$scope.hideGlobalLoader();

			// Check status
			if (data.Status === 'error') {

				alert(data.Message);
			
			} else if (data.Status === 'created') {
				
				// TODO: Animate transition between previously available element to
				// 'used' one, not sure whether we should use css3 animations or jQuery

				// Find machine by ID
				for (var machineIter = 0; machineIter < $scope.machines.length; machineIter++) {

					// Machine found condition
					if ($scope.machines[machineIter].Id === machineId) {
						
						// Refresh data of the new activation
						$scope.machines[machineIter].ActivationSecondsElapsed = 0;
						$scope.machines[machineIter].OccupiedByUserId = parseInt( $cookieStore.get('UserId') );
						$scope.machines[machineIter].ActivationId = data.Id;
						$scope.machines[machineIter].available = false;
						$scope.machines[machineIter].used = true;

						// Start timer for elapsed time
						$scope.machines[machineIter].activationInterval = 
							setInterval($scope.updateElapsedTime, 1000, machineIter);

						// Exit the machine finding for loop
						break;

					}

				}
			
			}

		})
		.error(function() {
			$scope.hideGlobalLoader();
			alert('Could not activate machine');
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
			console.log('Return to normal');
		});
	};

	$scope.deactivate = function(machine) {

		// Stop activation timer interval
		clearInterval(machine.activationInterval);

		// Show loader
		$scope.showGlobalLoader();

		$http({
			method: 'PUT', 
			url: '/api/activations', 
			params: {
				activation_id: machine.ActivationId,
				anticache: new Date().getTime()
			}
		})
		.success(function(data) {
			
			$scope.hideGlobalLoader();

			if (data.Status === 'ok') {

				// TODO: dynamicaly switch state of the previously
				// available item to 'used' using animation
				machine.used = false;
				machine.occupied = false;
				machine.available = true;
			
			} else if (data.Status === 'error') {
				alert(data.Message);
			} else {
				alert('Error while deactivating');
			}

		})
		.error(function() {
			$scope.hideGlobalLoader();
			alert('Failed to deactivate');
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

		console.log('show modal');

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
			if ($scope.machine.Status === MACHINE_STATUS_OCCUPIED) {
				
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
  	console.log('Log out and deactivate');
    $modalInstance.close('eh');
  };

  $scope.cancel = function () {
  	console.log('cancel');
    $modalInstance.dismiss('cancel');
  };

});

})(); // closure
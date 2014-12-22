(function(){

'use strict';

// Constants for machine states
var MACHINE_STATUS_AVAILABLE = 'free';
var MACHINE_STATUS_OCCUPIED = 'occupied';
var MACHINE_STATUS_USED = 'used'; // Used by the current user
var MACHINE_STATUS_UNAVAILABLE = 'unavailable';

// Our local app variable for the module
var app = angular.module('fabsmith.machines', ['ngRoute', 'timer', 'fabsmithFilters']);

app.config(['$routeProvider', function($routeProvider) {
	$routeProvider.when('/machines', {
		templateUrl: '/static/app/machines/machines.html',
		controller: 'MachinesCtrl'
	});
}]);

app.controller('MachinesCtrl', 
['$scope', '$http', '$location', '$route', '$cookieStore', 
function($scope, $http, $location, $route, $cookieStore) {

	// Show logged user name
	$scope.userFullName = $cookieStore.get('FirstName') + ' ' + $cookieStore.get('LastName');

	// Configure timer
	$scope.resetTimer = function() {
		$scope.$broadcast('timer-clear');
		$scope.$broadcast('timer-set-countdown', 120);
		$scope.$broadcast('timer-start');
	};
	$scope.resetTimer();

	// Initialize the machines array
	$scope.machines = [];

	$scope.onMachinesLoaded = function(machines) {

		// Add extra vars to machines
		for (var machinesIter = 0; machinesIter < machines.length; machinesIter++) {
			machines[machinesIter].available = $scope.isAvailable(machines[machinesIter]);
			machines[machinesIter].used = $scope.isUsed(machines[machinesIter]);
			machines[machinesIter].occupied = $scope.isOccupied(machines[machinesIter]);
			machines[machinesIter].unavailable = $scope.isUnavailable(machines[machinesIter]);
		}

		$scope.machines = machines;
		//console.log(machines);
	
	};

	// Get current user machines
	$http({
		method: 'GET',
		url: '/api/machines',
		params: {
			anticache: new Date().getTime()
		}
	})
	.success(function(data) {
		if (data.Status === 'error') {
			alert('msg: ' + data.Message);
		} else if (data.Machines.length <= 0) {
			alert('There are no machines available for you');
		} else if (data.Machines.length > 0) {
			//console.log(data.Machines);
			$scope.onMachinesLoaded(data.Machines);
		} else {
			alert('Error loading machines');
		}
	})
	.error(function() {
		alert('Error loading machines')
	});

	// Activate a machine by the currenty logged in user
	$scope.activate = function(machineId) {
		
		$scope.resetTimer();

		$http({
			method: 'POST',
			url: '/api/activations',
			params: {
				machine_id: machineId,
				anticache: new Date().getTime()
			}
		})
		.success(function(data) {
			
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
						$scope.machines[machineIter].activationInterval = setInterval(function() {
							$scope.machines[machineIter].ActivationSecondsElapsed++;
							$scope.$apply();
						}, 1000);

						// Exit the machine finding for loop
						break;

					}

				}
			
			}

		})
		.error(function() {
			alert('Could not activate machine');
		});

	};

	$scope.deactivate = function(machine) {
		
		$scope.resetTimer();

		if (!confirm('Make this machine available to other users')) {
			return;
		} 

		// Else continue
		// Stop activation timer interval
		clearInterval(machine.activationInterval);

		$http({
			method: 'PUT', 
			url: '/api/activations', 
			params: {
				activation_id: machine.ActivationId,
				anticache: new Date().getTime()
			}
		})
		.success(function(data) {
			
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
			alert('Failed to deactivate');
		});

	};

	$scope.isAvailable = function(machine) {
		return machine.Status === MACHINE_STATUS_AVAILABLE;
	}

	$scope.isOccupied = function(machine) {
		return machine.Status === MACHINE_STATUS_OCCUPIED;
	}

	$scope.isUsed = function(machine) {
		return machine.Status === MACHINE_STATUS_USED;
	}

	$scope.isUnavailable = function(machine) {
		return machine.Status === MACHINE_STATUS_UNAVAILABLE;
	}

	$scope.setAllStates = function(machine, trueOrFalse) {
		machine.available = trueOrFalse;
		machine.used = trueOrFalse;
		machine.occupied = trueOrFalse;
		machine.unavailable = trueOrFalse;
	}

}]);

app.directive('fsMachineItem', function() {
	return {
		templateUrl: 'static/app/machines/machine-item.html',
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
		templateUrl: 'static/app/machines/machine-body-available.html',
		restrict: 'E'
	}
});

app.directive('fsMachineBodyUsed', function() {
	return {
		templateUrl: 'static/app/machines/machine-body-used.html',
		restrict: 'E',
		controller: ['$scope', function($scope){

			$scope.timeElapsed = $scope.machine.ActivationSecondsElapsed;
			
			if ($scope.machine.Status == MACHINE_STATUS_USED) {
				$scope.machine.activationInterval = setInterval(function() {
					$scope.machine.ActivationSecondsElapsed++;
					$scope.$apply();
				}, 1000);
			}

		}]
	}
});

app.directive('fsMachineBodyOccupied', function() {
	return {
		templateUrl: 'static/app/machines/machine-body-occupied.html',
		restrict: 'E',
		controller: ['$scope', '$cookieStore', function($scope, $cookieStore){

			var user = {};
			user.Admin = $cookieStore.get('Admin');
			user.Staff = $cookieStore.get('Staff');
			user.Member = $cookieStore.get('Member');
			$scope.user = user;
			$scope.$apply();

		}]
	}
});

app.directive('fsMachineBodyUnavailable', function() {
	return {
		templateUrl: 'static/app/machines/machine-body-unavailable.html',
		restrict: 'E'
	}
});

})();
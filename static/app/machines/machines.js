(function(){

'use strict';

angular.module('fabsmith.machines', ['ngRoute'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/machines', {
    templateUrl: '/static/app/machines/machines.html',
    controller: 'MachinesCtrl'
  });
}])

.controller('MachinesListCtrl', ['$scope', '$http', function($scope, $http){
	$scope.test = 'Machines List'
}])

.controller('MachinesItemCtrl', ['$scope', '$http', function($scope, $http){
	$scope.test = 'Machines Item'
}])

.controller('MachinesCtrl', ['$scope', '$http', '$location', '$route', function($scope, $http, $location, $route) {

	// Constants
	var MACHINE_STATUS_AVAILABLE = 'free';
	var MACHINE_STATUS_OCCUPIED = 'occupied';
	var MACHINE_STATUS_USED = 'used'; // Used by the current user
	var MACHINE_STATUS_UNAVAILABLE = 'unavailable';

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
		console.log(machines);
	
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
				// 'used' one
				$route.reload();
			
			}
			//alert(data);
		})
		.error(function() {
			alert('Could not activate machine');
		});
	};

	$scope.deactivate = function(machine) {
		if (!confirm('Make this machine available to other users')) {
			return;
		}
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
				$route.reload();
			
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

	$scope.logout = function() {
		$http({
			method: 'GET',
			url: '/api/logout',
			params: {
				anticache: new Date().getTime()
			}
		})
		.success(function() {
			$location.path('/');
		})
		.error(function() {
			alert('Failed to log out. Probably server down.');
		});
	};
}])

.directive('fsMachineItem', function() {
	return {
    templateUrl: 'static/app/machines/machine-item.html',
    restrict: 'E'
  };
})

.directive('fsMachineBodyAvailable', function() {
	return {
		templateUrl: 'static/app/machines/machine-body-available.html',
		restrict: 'E'
	}
})

.directive('fsMachineBodyUsed', function() {
	return {
		templateUrl: 'static/app/machines/machine-body-used.html',
		restrict: 'E'
	}
})

.directive('fsMachineBodyOccupied', function() {
	return {
		templateUrl: 'static/app/machines/machine-body-occupied.html',
		restrict: 'E'
	}
})

.directive('fsMachineBodyUnavailable', function() {
	return {
		templateUrl: 'static/app/machines/machine-body-unavailable.html',
		restrict: 'E'
	}
});

})();
(function(){

'use strict';

angular.module('fabsmith.machines', ['ngRoute'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/machines', {
    templateUrl: '/static/app/machines/machines.html',
    controller: 'MachinesCtrl'
  });
}])

.controller('ActivationsCtrl', ['$scope', '$http', '$location', function($scope, $http, $location){
	$scope.activations = [];
	// Get current user activations on load
	$http({
		method: 'GET',
		url: '/api/activations',
		params: {
			anticache: new Date().getTime()
		}
	})
	.success(function(data) {
		if (data.length && data[0].Id) {
			$scope.activations = data;
		} else if (data.Status == 'error') {
			//alert(data.Message);
		} else {
			alert('Failed to get current activations');
		}
	})
	.error(function() {
		alert('Error getting activations');
	});

	$scope.deactivate = function(activation) {
		if (!confirm('Make this machine available to other users')) {
			return;
		}
		$http({
			method: 'PUT', 
			url: '/api/activations', 
			params: {
				activation_id: activation.Id,
				anticache: new Date().getTime()
			}
		})
		.success(function(data) {
			if (data.Status === 'ok') {
				$scope.activations.splice($scope.activations.indexOf(activation), 1);
				$location.path('/logout');
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

	$scope.test = 'Activations';
}])

.controller('MachinesListCtrl', ['$scope', '$http', function($scope, $http){
	$scope.test = 'Machines List'
}])

.controller('MachinesItemCtrl', ['$scope', '$http', function($scope, $http){
	$scope.test = 'Machines Item'
}])

.controller('MachinesCtrl', ['$scope', '$http', '$location', function($scope, $http, $location) {
	$scope.test = 'Machines'
	$scope.machines = [];
	// Load machines
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
			$scope.machines = data.Machines;
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
				// Activation successful - redirect to logout countdown
				//alert(data.Id);
				$location.path('/logout');
			}
			//alert(data);
		})
		.error(function() {
			alert('Could not activate machine');
		});
	};

	$scope.isFree = function(machine) {
		return machine.Status === 'free';
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
}]);

})();
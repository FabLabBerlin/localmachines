(function(){

'use strict';
var app = angular.module('fabsmith.admin.user', ['ngRoute', 'ngCookies', 'fabsmith.admin.randomtoken']);

app.config(['$routeProvider', function($routeProvider) {
	$routeProvider.when('/users/:userId', {
		templateUrl: 'ng-modules/user/user.html',
		controller: 'UserCtrl'
	});
}]); // app.config

app.controller('UserCtrl', ['$scope', '$routeParams', '$http', '$location', 'randomToken', function($scope, $routeParams, $http, $location, randomToken) {
	
	$('.datepicker').pickadate();

	$scope.user = {
		Id: $routeParams.userId
	};
	$scope.userMachines = [];
	$scope.userMemberships = [];

	$http({
		method: 'GET',
		url: '/api/users/' + $scope.user.Id,
		params: {
			anticache: new Date().getTime()
		}
	})
	.success(function(data) {
		console.log('Got user');
		console.log(data);
		$scope.user = data;
	})
	.error(function(data, status) {
		console.log('Could not get user');
		console.log('Data: ' + data);
		console.log('Status code: ' + status);
	});

	$http({
		method: 'GET',
		url: '/api/users/' + $scope.user.Id + '/machines',
		params: {
			anticache: new Date().getTime()
		}
	})
	.success(function(data) {
		console.log('Got user machines');
		console.log(data);
		$scope.userMachines = data;
	})
	.error(function(data, status) {
		console.log('Could not get user machines');
		console.log('Data: ' + data);
		console.log('Status code: ' + status);
	});

	function formatDate(d) {
		console.log('formatDate: d: ', d);
		var mm = (d.getMonth() + 1);
		if (mm < 10) {
			mm = '0' + mm;
		}
		var dd = d.getDate();
		if (dd < 10) {
			dd = '0' + dd;
		}
		return d.getFullYear() + '-' + mm + '-' + dd;
	}

	$http({
		method: 'GET',
		url: '/api/memberships',
		params: {
			anticache: new Date().getTime()
		}
	})
	.success(function(data) {
		$scope.memberships = data;
		$scope.membershipsById = {};
		_.each($scope.memberships, function(m) {
			$scope.membershipsById[m.Id] = m;
		});

		$http({
			method: 'GET',
			url: '/api/users/' + $scope.user.Id + '/memberships',
			params: {
				anticache: new Date().getTime()
			}
		})
		.success(function(data) {
			console.log('Got user memberships');
			console.log(data);
			$scope.userMemberships = _.map(data, function(userMembership) {
				userMembership.StartDate = new Date(Date.parse(userMembership.StartDate));
				_.merge(userMembership, {
					EndDate: new Date(userMembership.StartDate)
				});
				console.log('userMembership.Id: ', userMembership.Id);
				var membership = $scope.membershipsById[userMembership.MembershipId];
				console.log('membership: ', membership);
				userMembership.EndDate.setDate(userMembership.StartDate.getDate() + membership.Duration); 
				userMembership.StartDate = formatDate(userMembership.StartDate);
				userMembership.EndDate = formatDate(userMembership.EndDate);
				return userMembership;
			});
		})
		.error(function(data, status) {
			console.log('Could not get user memberships');
			console.log('Data: ' + data);
			console.log('Status code: ' + status);
		});
	})
	.error(function(data, status) {
		console.log('Could not get memberships');
		console.log('Data: ' + data);
		console.log('Status code: ' + status);
	});

	$scope.addUserMembership = function() {
		var startDate = $('#adm-add-user-membership-start-date').val();
		if (!startDate) {
			toastr.error('Please select a Start Date');
			return;
		}
		startDate = new Date(startDate);
		startDate = formatDate(startDate);
		var userMembershipId = $('#user-select-membership').val();
		if (!userMembershipId) {
			toastr.error('Please select a Membership');
			return;
		}
		$http({
			method: 'POST',
			url: '/api/users/' + $scope.user.Id + '/memberships',
			data: {
				StartDate: startDate,
				UserMembershipId: userMembershipId
			}
		})
		.success(function() {
			toastr.info('Successfully created user membership');
			window.location.reload(true);
		})
		.error(function() {
			toastr.error('Error while trying to create new User Membership');
		});
	};

	$scope.cancel = function() {
		if (confirm('All changes will be discarded, click ok to continue.')) {
			$location.path('/users');
		}
	};

	$scope.deleteUser = function() {
		var email = prompt("Do you really want to delete this user? Please enter user's E-Mail address to continue");
		if (email === $scope.user.Email) {
			$http({
				method: 'DELETE',
				url: '/api/users/' + $scope.user.Id
			})
			.success(function(data) {
				toastr.info('User deleted');
				$location.path('/users');
			})
			.error(function() {
				toastr.error('Error while trying to delete user');
			});
		} else {
			toastr.warning('Delete User canceled.');
		}
	};

	$scope.deleteUserMembership = function(userMembershipId) {
		var token = randomToken.generate();
		vex.dialog.prompt({
			message: 'Enter <span class="delete-prompt-token">' +
			token + '</span> to delete',
			placeholder: 'Token',
			callback: function(value) {
				if (!value) {
					toastr.error('No token');
				} else if (value === token) {
					$http({
						method: 'DELETE',
						url: '/api/users/' + $scope.user.Id + '/memberships/' + userMembershipId
					})
					.success(function(data) {
						toastr.info('Successfully deleted user membership');
						window.location.reload(true);
					})
					.error(function() {
						toastr.error('Error while trying to delete user membership');
					});
				} else {
					toastr.error('Wrong token');
				}
			}
		});
	};

	$scope.saveUser = function() {
		console.log('user model:', $scope.user);
		$http({
			method: 'PUT',
			url: '/api/users/' + $scope.user.Id,
			headers: {'Content-Type': 'application/json' },
			data: {
				User: $scope.user,
				UserRoles: $scope.userRoles
			},
			transformRequest: function(data) {
				console.log('data to send:', data);
				return JSON.stringify(data);
			}
		})
		.success(function() {
			toastr.info('Changes saved.');
		})
		.error(function() {
			toastr.error('Error while trying to save changes');
		});
	};

	$scope.updateUserRoles = function(value) {
		console.log('value:',value);
	};
}]); // app.controller

})(); // closure
(function(){

'use strict';

var app = angular.module('fabsmith.admin.membership', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/membership/:membershipId', {
    templateUrl: 'ng-modules/membership/membership.html',
    controller: 'MembershipCtrl'
  });
}]); // app.config

app.controller('MembershipCtrl', 
 ['$scope', '$http', '$location', '$filter', '$routeParams',
 function($scope, $http, $location, $filter, $routeParams) {

  $scope.machines = [];
  $scope.membership = {
    Id: $routeParams.membershipId
  };

  // Load machines first
  $http({
    method: 'GET',
    url: '/api/machines',
    params: {
      anticache: new Date().getTime()
    }
  })
  .success(function(machines) {
    $scope.machines = machines;
    $scope.loadMembership();
  })
  .error(function() {
    toastr.error('Failed to load machines');
  });

  // Load membership
  $scope.loadMembership = function() {
    $http({
      method: 'GET',
      url: '/api/memberships/' + $scope.membership.Id,
      params: {
        anticache: new Date().getTime()
      }
    })
    .success(function(membershipModel) {
      $scope.membership = membershipModel;
      
      // Parse affected machines JSON as it is passed here as string
      $scope.membership.AffectedMachines = 
       JSON.parse($scope.membership.AffectedMachines);

      // Search for machines with the same IDs as the AffectedMachines
      // and set them as checked
      for (var i = 0; i < $scope.membership.AffectedMachines.length; i++) {
        for (var j = 0; j < $scope.machines.length; j++) {
          if (parseInt($scope.machines[j].Id) === 
              parseInt($scope.membership.AffectedMachines[i])) {
            $scope.machines[j].Checked = true;
          }
        }
      }

      // Filter as currency
      $scope.membership.Price = 
       $filter('currency')($scope.membership.Price,'',2);
    })
    .error(function() {
      toastr.error('Failed to get membership');
    });
  };

  $scope.updateMembership = function() {

    // Add the machine.Checked's into membership.AffectedMachines array
    var affectedMachines = [];
    for (var i = 0; i < $scope.machines.length; i++) {
      if ($scope.machines[i].Checked) {
        affectedMachines.push($scope.machines[i].Id);
      }
    }
    affectedMachines = JSON.stringify(affectedMachines);
    
    // Make a clone of the model to feel safe
    var membership = _.clone($scope.membership);

    membership.AffectedMachines = affectedMachines;
    membership.Price = parseFloat(membership.Price);
    membership.Duration = parseInt(membership.Duration);
    membership.MachinePriceDeduction = 
     parseInt(membership.MachinePriceDeduction);

    $http({
      method: 'PUT',
      url: '/api/memberships/' + membership.Id,
      headers: {'Content-Type': 'application/json'},
      data: membership,
      transformRequest: function(data) {
        console.log('Membership data to send:', data);
        return JSON.stringify(data);
      },
      params: {
        anticache: new Date().getTime()
      }
    })
    .success(function() {
      toastr.success('Membership updated');
    })
    .error(function() {
      toastr.error('Failed to update membership');
    });
  }; // updateMembership()

  $scope.deleteMembershipPrompt = function() {

    // You have to add the <random-token></random-token> directive somewhere
    // in HTML in order to make this work
    $scope.generateRandomToken();

    vex.dialog.prompt({
      // Unfortunately it is not possible to parse directives inside
      // vex messages, so we just get the random token
      message: 'Enter <span class="delete-prompt-token">' + 
       $scope.randomToken + '</span> to delete',
      placeholder: 'Token',
      callback: $scope.deleteMembershipPromptCallback
    });
  };

  $scope.deleteMembershipPromptCallback = function(value) {
    if (value) {    
      if (value === $scope.randomToken) {
        $scope.deleteMembership();
      } else {
        toastr.error('Wrong token');
      }
    } else if (value !== false) {
      toastr.error('No token');
    }
  };

  $scope.deleteMembership = function() {
    $http({
      method: 'DELETE',
      url: '/api/memberships/' + $scope.membership.Id,
      params: {
        anticache: new Date().getTime()
      }
    })
    .success(function() {
      toastr.success('Membership deleted');
      $location.path('/memberships');
    })
    .error(function() {
      toastr.error('Failed to delete membership');
    });
  };

}]); // app.controller

})(); // closure
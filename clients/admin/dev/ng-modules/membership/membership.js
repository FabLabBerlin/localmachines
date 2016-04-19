(function(){

'use strict';

var app = angular.module('fabsmith.admin.membership', ['ngRoute', 'ngCookies', 'fabsmith.admin.randomtoken']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/membership/:membershipId', {
    templateUrl: 'ng-modules/membership/membership.html',
    controller: 'MembershipCtrl'
  });
}]); // app.config

app.controller('MembershipCtrl', 
 ['$scope', '$http', '$location', '$filter', '$routeParams', 'randomToken', 'api',
 function($scope, $http, $location, $filter, $routeParams, randomToken, api) {

  $scope.machines = [];
  $scope.membership = {
    Id: $routeParams.membershipId,
    AutoExtend: true,     // default
    AutoExtendDuration: 1 // values
  };

  // Load machines first
  api.loadMachines(function(resp) {
    $scope.machines = resp.machines;
    $scope.loadMembership();
  });

  // Load membership
  $scope.loadMembership = function() {
    $http({
      method: 'GET',
      url: '/api/memberships/' + $scope.membership.Id,
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(membershipModel) {
      $scope.membership = membershipModel;
      
      if ($scope.membership.AffectedMachines !== '') {
        // Parse affected machines JSON as it is passed here as string
        $scope.membership.AffectedMachines = 
          JSON.parse($scope.membership.AffectedMachines);
      }

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
      $scope.membership.MonthlyPrice = 
       $filter('currency')($scope.membership.MonthlyPrice,'',2);
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
    membership.MonthlyPrice = parseFloat(membership.MonthlyPrice);
    membership.DurationMonths = parseInt(membership.DurationMonths);
    membership.MachinePriceDeduction = 
     parseInt(membership.MachinePriceDeduction);
    membership.AutoExtendDurationMonths = 
     parseInt(membership.AutoExtendDurationMonths);

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
        ac: new Date().getTime()
      }
    })
    .success(function() {
      toastr.success('Membership updated');
    })
    .error(function() {
      toastr.error('Failed to update membership');
    });
  }; // updateMembership()

  $scope.setArchived = function(archived) {
    var action = archived ? 'archived' : 'unarchived';

    $http({
      method: 'POST',
      url: '/api/memberships/' + $scope.membership.Id + '/set_archived',
      params: {
        archived: archived
      }
    })
    .success(function(data) {
      toastr.info('Successfully ' + action + ' membership');
      $scope.loadMembership();
    })
    .error(function() {
      toastr.error('Failed to ' + action + ' membership');
    });
  }; // archive()

}]); // app.controller

})(); // closure
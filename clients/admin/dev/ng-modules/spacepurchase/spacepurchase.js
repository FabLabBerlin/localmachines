(function(){

'use strict';
var app = angular.module('fabsmith.admin.space.purchase', ['ngRoute', 'ngCookies', 'fabsmith.admin.randomtoken']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/space_purchases/:id', {
    templateUrl: 'ng-modules/spacepurchase/spacepurchase.html',
    controller: 'SpacePurchaseCtrl'
  });
}]); // app.config

app.controller('SpacePurchaseCtrl',
 ['$scope', '$routeParams', '$http', '$location', 'randomToken',
 function($scope, $routeParams, $http, $location, randomToken) {

  $scope.purchases = [];
  $scope.spacePurchase = {
    Id: $routeParams.id
  };
  $scope.users = [];
  $scope.usersById = {};

  function loadSpaces() {
    $http({
      method: 'GET',
      url: '/api/spaces',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      $scope.spaces = _.sortBy(data, function(space) {
        return space.Product.Name;
      });
    })
    .error(function() {
      toastr.error('Failed to get spaces');
    });
  }

  function loadSpacePurchase() {
    $http({
      method: 'GET',
      url: '/api/space_purchases/' + $scope.spacePurchase.Id,
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(sp) {
      $scope.spacePurchase = sp;
      sp.TimeStartLocal = moment(sp.TimeStart).tz('Europe/Berlin').format('YYYY-MM-DD HH:mm');
      sp.TimeEndLocal = moment(sp.TimeEnd).tz('Europe/Berlin').format('YYYY-MM-DD HH:mm');
    })
    .error(function(data, status) {
      toastr.error('Failed to load user data');
    });
  }

  function loadUsers() {
    $http({
      method: 'GET',
      url: '/api/users',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      $scope.users = _.sortBy(data, function(user) {
        return user.FirstName + ' ' + user.LastName;
      });
      _.each($scope.users, function(user) {
        $scope.usersById[user.Id] = user;
      });
    })
    .error(function() {
      toastr.error('Failed to get reservations');
    });
  }

  $scope.updateSpacePurchase = function() {
    $http({
      method: 'PUT',
      url: '/api/space_purchases/' + $scope.spacePurchase.Id,
      headers: {'Content-Type': 'application/json' },
      data: $scope.spacePurchase,
      transformRequest: function(data) {
        return JSON.stringify(data);
      },
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      toastr.success('Update successful');
    })
    .error(function(data) {
      console.log(data);
      toastr.error('Failed to update');
    });
  };

  $scope.save = function() {
    $http({
      method: 'PUT',
      url: '/api/space_purchases/' + $scope.spacePurchase.Id,
      headers: {'Content-Type': 'application/json' },
      data: $scope.spacePurchase,
      transformRequest: function(data) {
        var transformed = _.extend({}, data);
        transformed.SpaceId = parseInt(data.SpaceId);
        transformed.TimeStart = moment.tz(data.TimeStartLocal, 'Europe/Berlin').toDate();
        transformed.TimeEnd = moment.tz(data.TimeEndLocal, 'Europe/Berlin').toDate();
        transformed.UserId = parseInt(data.UserId);
        transformed.Quantity = parseInt(data.Quantity);
        console.log('transformed:', transformed);
        return JSON.stringify(transformed);
      },
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function() {
      toastr.success('Space purchase updated');
    })
    .error(function(data) {
      toastr.error('Error while trying to save changes');
    });
  };

  $scope.remove = function() {
    var token = randomToken.generate();
    vex.dialog.prompt({
      message: 'Enter <span class="delete-prompt-token">' +
       token + '</span> to delete',
      placeholder: 'Token',
      callback: function(value) {
        if (value) {
          if (value === token) {
            $http({
              method: 'DELETE',
              url: '/api/space_purchases/' + $scope.spacePurchase.Id,
              params: {
                ac: new Date().getTime()
              }
            })
            .success(function() {
              toastr.success('Space purchase deleted');
              $location.path('/spaces');
            })
            .error(function() {
              toastr.error('Error while trying to delete Space purchase');
            });
          } else {
            toastr.error('Wrong token');
          }
        } else {
          toastr.warning('Deletion canceled');
        }
      }
    });
  };

  loadSpaces();
  loadSpacePurchase();
  loadUsers();

}]); // app.controller

})(); // closure

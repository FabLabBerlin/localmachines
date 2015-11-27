(function(){

'use strict';

var mod = angular.module("fabsmith.admin.api", []);

mod.service('api', function($http) {
  this.loadSpaces = function(cb) {
    $http({
      method: 'GET',
      url: '/api/products',
      params: {
        ac: new Date().getTime(),
        type: 'space'
      }
    })
    .success(function(data) {
      var spaces = _.sortBy(data, function(space) {
        return space.Product.Name;
      });
      var spacesById = {};
      _.each(spaces, function(space) {
        spacesById[space.Product.Id] = space;
      });
      if (cb) {
        cb({
          spaces: spaces,
          spacesById: spacesById
        });
      }
    })
    .error(function() {
      toastr.error('Failed to get spaces');
    });
  };

  this.loadUsers = function(cb) {
    $http({
      method: 'GET',
      url: '/api/users',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      var users = _.sortBy(data, function(user) {
        return user.FirstName + ' ' + user.LastName;
      });
      var usersById = {};
      _.each(users, function(user) {
        usersById[user.Id] = user;
      });
      if (cb) {
        cb({
          users: users,
          usersById: usersById
        });
      }
    })
    .error(function() {
      toastr.error('Failed to get reservations');
    });
  };

  return this;
});

})(); // closure

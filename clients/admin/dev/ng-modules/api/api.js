(function(){

'use strict';

var mod = angular.module("fabsmith.admin.api", []);

mod.service('api', function($http) {
  // Public Methods

  this.loadMachines = function(cb) {
    $http({
      method: 'GET',
      url: '/api/machines',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(machines) {
      var machinesById = {};
      _.each(machines, function(machine) {
        machinesById[machine.Id] = machine;
      });
      if (cb) {
        cb({
          machines: machines,
          machinesById: machinesById
        });
      }
    })
    .error(function() {
      toastr.error('Failed to get machines');
    });
  };

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

  this.loadSpacePurchase = function(id, cb) {
    $http({
      method: 'GET',
      url: '/api/purchases/' + id,
      params: {
        ac: new Date().getTime(),
        type: 'space'
      }
    })
    .success(function(sp) {
      generateStartEndDateTimesLocal(sp);
      if (cb) {
        cb(sp);
      }
    })
    .error(function(data, status) {
      toastr.error('Failed to load space purchase data');
    });
  };

  this.loadTutoringPurchase = function(id, cb) {
    $http({
      method: 'GET',
      url: '/api/purchases/' + id,
      params: {
        ac: new Date().getTime(),
        type: 'tutor'
      }
    })
    .success(function(sp) {
      generateStartEndDateTimesLocal(sp);
      if (cb) {
        cb(sp);
      }
    })
    .error(function(data, status) {
      toastr.error('Failed to load tutoring purchase data');
    });
  };

  this.loadTutors = function(cb) {
    $http({
      method: 'GET',
      url: '/api/tutoring/tutors',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(tutorList) {
      var tutorsById = {};
      _.each(tutorList.Data, function(tutor) {
        tutorsById[tutor.Id] = tutor;
      });
      if (cb) {
        cb({
          tutors: tutorList.Data,
          tutorsById: tutorsById
        });
      }
    })
    .error(function() {
      toastr.error('Failed to load tutor list');
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

  this.purchase = {
    calculateQuantity: function(purchase) {
      console.log('$scope.timeChange()');
      this.parseInputTimes(purchase);
      var start = moment(purchase.TimeStart);
      var end = moment(purchase.TimeEnd);
      var duration = end.unix() - start.unix();
      console.log('duration=', duration);
      var quantity;
      switch (purchase.PriceUnit) {
      case 'minute':
        quantity = duration / 60;
        break;
      case 'hour':
        quantity = duration / 3600;
        break;
      case 'day':
        quantity = duration / 24 / 3600;
        break;
      default:
        return;
      }
      purchase.Quantity = quantity;
    },

    calculateTotalPrice: function(purchase) {
      var totalPrice = purchase.Quantity * purchase.PricePerUnit;
      purchase.TotalPrice = totalPrice.toFixed(2);
    },

    parseInputTimes: function(purchase) {
      var p = purchase;
      p.TimeStart = moment.tz(p.DateStartLocal + ' ' + p.TimeStartLocal, 'Europe/Berlin').toDate();
      p.TimeEnd = moment.tz(p.DateEndLocal + ' ' + p.TimeEndLocal, 'Europe/Berlin').toDate();
    }
  };


  // Private methods

  function generateStartEndDateTimesLocal(purchase) {
      var start = moment(purchase.TimeStart).tz('Europe/Berlin');
      var end = moment(purchase.TimeEnd).tz('Europe/Berlin');
      purchase.DateStartLocal = start.format('YYYY-MM-DD');
      purchase.DateEndLocal = end.format('YYYY-MM-DD');
      purchase.TimeStartLocal = start.format('HH:mm');
      purchase.TimeEndLocal = end.format('HH:mm');
  }

  return this;
});

})(); // closure

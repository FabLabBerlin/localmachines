(function(){

'use strict';

var mod = angular.module("fabsmith.admin.api", []);

mod.service('api',
 ['$cookies', '$http', 'randomToken',
 function($cookies, $http, randomToken) {
  // Public Methods

  this.prompt = function(msg, cb) {
    var expectedToken = randomToken.generate();
    vex.dialog.prompt({
      message: msg + ' Enter <span class="delete-prompt-token">' + 
       expectedToken + '</span> to continue.',
      placeholder: 'Token',
      callback: function(value) {
        if (value) {
          if (value === expectedToken) {
            cb();
          } else {
            toastr.error('Wrong token');
          }
        } else if (value !== false) {
          toastr.error('No token');
        }
      }
    });
  };

  this.loadMachines = function(cb) {
    $http({
      method: 'GET',
      url: '/api/machines',
      params: {
        ac: new Date().getTime(),
        location: $cookies.locationId
      }
    })
    .success(function(machines) {
      var machinesById = {};
      machines = _.sortBy(machines, function(m) {
        return m.Name;
      });
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
        location: $cookies.locationId,
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
      url: '/api/products',
      params: {
        location: $cookies.locationId,
        ac: new Date().getTime(),
        type: 'tutor'
      }
    })
    .success(function(tutorList) {
      var tutorsById = {};
      _.each(tutorList, function(tutor) {
        tutorsById[tutor.Product.Id] = tutor;
      });
      if (cb) {
        cb({
          tutors: tutorList,
          tutorsById: tutorsById
        });
      }
    })
    .error(function() {
      toastr.error('Failed to load tutor list');
    });
  };

  this.archiveProduct = function(id, success, error) {
    $http({
      method: 'PUT',
      url: '/api/products/' + id + '/archive',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function() {
      if (success) {
        success();
      }
    })
    .error(function() {
      if (error) {
        error();
      }
    });
  };

  this.archivePurchase = function(id, success, error) {
    $http({
      method: 'PUT',
      url: '/api/purchases/' + id + '/archive',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function() {
      if (success) {
        success();
      }
    })
    .error(function() {
      if (error) {
        error();
      }
    });
  };

  this.loadUsers = function(cb) {
    $http({
      method: 'GET',
      url: '/api/users',
      params: {
        ac: new Date().getTime(),
        location: $cookies.locationId
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

  this.loadPurchases = function(success, error) {
    $http({
      method: 'GET',
      url: '/api/purchases',
      params: {
        location: $cookies.locationId,
        ac: new Date().getTime(),
        type: 'tutor'
      }
    })
    .success(function(purchaseList) {
      if (success) {
        success(purchaseList);
      }
    })
    .error(function() {
      if (error) {
        error();
      }
    });
  };

  this.purchase = {
    calculateQuantity: function(purchase) {
      console.log('$scope.timeChange()');
      this.parseInputTimes(purchase);
      console.log('calculateQuantity: purchase.TimeStart=', purchase.TimeStart);
      console.log('calculateQuantity: purchase.TimeEnd=', purchase.TimeEnd);
      console.log('calculateQuantity: purchase.TimeEndPlanned=', purchase.TimeEndPlanned);
      var start = moment(purchase.TimeStart);
      var end = moment(purchase.TimeEnd);
      var endPlanned = moment(purchase.TimeEndPlanned);
      if (start.unix() > 0) {
        console.log('s>0');
        var duration;
        if (end.unix() > 0) {
          console.log('e>0');
          duration = end.unix() - start.unix();
        } else if (endPlanned.unix() > 0) {
          console.log('ep>0');
          duration = endPlanned.unix() - start.unix();
        } else {
          return;
        }
        console.log('duration=', duration);
        var quantity;
        console.log('purchase.PriceUnit=', purchase.PriceUnit);
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
      }
    },

    calculateTotalPrice: function(purchase) {
      var totalPrice = purchase.Quantity * purchase.PricePerUnit;
      purchase.TotalPrice = totalPrice.toFixed(2);
    },

    parseInputTimes: function(purchase) {
      var p = purchase;
      p.TimeStart = moment.tz(p.DateStartLocal + ' ' + p.TimeStartLocal, 'Europe/Berlin').toDate();
      p.TimeEnd = moment.tz(p.DateEndLocal + ' ' + p.TimeEndLocal, 'Europe/Berlin').toDate();
      p.TimeEndPlanned = moment.tz(p.DateEndPlannedLocal + ' ' + p.TimeEndPlannedLocal, 'Europe/Berlin').toDate();
    }
  };

  this.toMoment = function(goDateTime, tz) {
    if (goDateTime && moment(goDateTime).unix() > 0) {
      return moment(goDateTime).tz(tz || 'Europe/Berlin');
    }
  };

  // Private methods

  function generateStartEndDateTimesLocal(purchase) {
      var start = moment(purchase.TimeStart).tz('Europe/Berlin');
      var end = moment(purchase.TimeEnd).tz('Europe/Berlin');
      var endPlanned = moment(purchase.TimeEndPlanned).tz('Europe/Berlin');
      purchase.DateStartLocal = start.format('YYYY-MM-DD');
      purchase.DateEndLocal = end.format('YYYY-MM-DD');
      purchase.DateEndPlannedLocal = endPlanned.format('YYYY-MM-DD');
      purchase.TimeStartLocal = start.format('HH:mm');
      purchase.TimeEndLocal = end.format('HH:mm');
      purchase.TimeEndPlannedLocal = endPlanned.format('HH:mm');
  }

  return this;
}]);

})(); // closure

(function(){

'use strict';

var app = angular.module('fabsmith.admin.tutoring', ['ngRoute', 'ngCookies', 'fabsmith.admin.api']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/tutoring', {
    templateUrl: 'ng-modules/tutoring/tutoring.html',
    controller: 'TutoringCtrl'
  });
}]); // app.config

app.controller('TutoringCtrl', ['$scope', '$http', '$location', 'api', 'randomToken',
  function($scope, $http, $location, api, randomToken) {

  $scope.machines = [];
  $scope.tutors = [];

  // Load global settings for the VAT and currency
  $scope.loadSettings = function() {
    $http({
      method: 'GET',
      url: '/api/settings',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(settings) {
      $scope.settings = {
        Currency: {},
        VAT: {}
      };
      _.each(settings, function(setting) {
        $scope.settings[setting.Name] = setting;
      });
    })
    .error(function() {
      toastr.error('Failed to get global config');
    });
  };

  $scope.showTutorSkills = function() {
    if ($scope.machines.length && $scope.tutors.length) {
      // Translate MachineSkills machine ID's to machine names for each tutor
      _.each($scope.tutors, function(tutor) {
        var machineSkills = [];
        if (tutor.Product.MachineSkills) {
          var tmp = tutor.Product.MachineSkills;
          tmp = tmp.slice(1, tmp.length - 1);
          _.each(tmp.split(','), function(idString) {
            var id = parseInt(idString);
            var machine = $scope.machinesById[id];
            if (machine) {
              machineSkills.push(machine);
            }
          });
        }
        tutor.MachineSkills = machineSkills;

      });
    }
  };

  $scope.getMachineNameById = function(machineId) {
    var machineName = '';
    _.each($scope.machines, function(machine) {
      if (parseInt(machine.Id) === parseInt(machineId)) {
        machineName = machine.Name;
      }
    });
    return machineName;
  };

  $scope.loadPurchases = function() {
    api.loadPurchases(function(purchaseList) {
      $scope.purchases = _.sortBy(purchaseList.Data, function(purchase) {
        return purchase.Name;
      });
      $scope.purchases = _.map($scope.purchases, function(p) {
        var tutor = $scope.tutorsById[p.ProductId];
        if (tutor) {
          p.Product = tutor.Product;
        }
        p.User = $scope.usersById[p.UserId];

        var timeCreated = api.toMoment(p.Created);
        var timeStart = api.toMoment(p.TimeStart);
        var timeEnd = api.toMoment(p.TimeEnd);

        p.Created = timeCreated && timeCreated.format('D MMM YYYY');
        p.TimeStart = timeStart && timeStart.format('D MMM YYYY HH:mm');
        p.TimeEnd = timeEnd && timeEnd.format('D MMM YYYY HH:mm');
        
        // Show reserved time
        if (timeStart && timeEnd) {
          p.TimeReserved = moment.duration(timeEnd.diff(timeStart)).format('h[h] m[m]');
        }

        // Show timer time
        var timerDuration;
        if (p.PriceUnit === 'day') {
          timerDuration = moment.duration(p.Quantity, 'days');
        } else if (p.PriceUnit === 'hour') {
          timerDuration = moment.duration(p.Quantity, 'hours');
        } else if (p.PriceUnit === 'minute') {
          timerDuration = moment.duration(p.Quantity, 'minutes');
        }

        if (timerDuration.asHours() < 1) {
          p.TimeTimed = timerDuration.format('m[m] s[s]');
        } else {
          p.TimeTimed = timerDuration.format('h[h] m[m]');
        }

        if (p.Quantity && p.PricePerUnit) {
          p.TotalPrice = p.Quantity * p.PricePerUnit;
        }
        return p;
      });

      $scope.showTutorSkills();
    }, function() {
      toastr.error('Failed to load purchases');
    });
  };

  $scope.addTutor = function() {
    $http({
      method: 'POST',
      url: '/api/products',
      params: {
        name: name,
        ac: new Date().getTime(),
        type: 'tutor'
      }
    })
    .success(function(data) {
      $scope.editTutor(data.Product.Id);
    })
    .error(function() {
      toastr.error('Failed to create product');
    });
  };

  $scope.editTutor = function(id) {
    $location.path('/tutoring/tutors/' + id);
  };

  $scope.archiveTutorPrompt = function(tutorId) {
    var token = randomToken.generate();
    vex.dialog.prompt({
      message: 'Enter <span class="delete-prompt-token">' + 
       token + '</span> to archive tutor',
      placeholder: 'Token',
      callback: $scope.archiveTutorPromptCallback.bind(this, token, tutorId)
    });
  };

  $scope.archiveTutorPromptCallback = function(expectedToken, tutorId, value) {
    if (value) {    
      if (value === expectedToken) {
        $scope.archiveTutor(tutorId);
      } else {
        toastr.error('Wrong token');
      }
    } else if (value !== false) {
      toastr.error('No token');
    }
  };

  $scope.archiveTutor = function(tutorId) {
    api.archiveProduct(
      tutorId,
      function(){
        toastr.success("Tutor has been archived");
        api.loadTutors(function(tutorData) {
          $scope.tutors = tutorData.tutors;
          $scope.tutorsById = tutorData.tutorsById;
          $scope.showTutorSkills();
        });
      }, 
      function() {
        toastr.error("Failed to archive tutor");
      });
  };

  $scope.addPurchase = function() {
    $http({
      method: 'POST',
      url: '/api/purchases',
      params: {
        ac: new Date().getTime(),
        type: 'tutor'
      }
    })
    .success(function(tutoringPurchase) {
      $location.path('/tutoring/purchases/' + tutoringPurchase.Id);
    })
    .error(function() {
      toastr.error('Failed to create tutoring purchase');
    });
  };

  $scope.editPurchase = function(id) {
    $location.path('/tutoring/purchases/' + id);
  };

  $scope.loadSettings();
  api.loadMachines(function(machineData) {
    $scope.machines = machineData.machines;
    $scope.machinesById = machineData.machinesById;
    api.loadTutors(function(tutorData) {
      $scope.tutors = tutorData.tutors;
      $scope.tutorsById = tutorData.tutorsById;
      api.loadUsers(function(userData) {
        $scope.users = userData.users;
        $scope.usersById = userData.usersById;
        $scope.loadPurchases();
      });
    });
  });

  $scope.archivePurchasePrompt = function(purchaseId) {
    var token = randomToken.generate();
    vex.dialog.prompt({
      message: 'Enter <span class="delete-prompt-token">' + 
       token + '</span> to archive tutor',
      placeholder: 'Token',
      callback: $scope.archivePurchasePromptCallback.bind(this, token, purchaseId)
    });
  };

  $scope.archivePurchasePromptCallback = function(expectedToken, purchaseId, value) {
    if (value) {    
      if (value === expectedToken) {
        $scope.archivePurchase(purchaseId);
      } else {
        toastr.error('Wrong token');
      }
    } else if (value !== false) {
      toastr.error('No token');
    }
  };

  $scope.archivePurchase = function(purchaseId) {
    api.archivePurchase(
      purchaseId,
      function(){
        toastr.success("Purchase has been archived");
        $scope.loadPurchases();
      }, 
      function() {
        toastr.error("Failed to archive purchase");
      });
  };

}]); // app.controller

})(); // closure
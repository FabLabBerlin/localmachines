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
          console.log('tmp:', tmp);
          tmp = tmp.slice(1, tmp.length - 1);
          console.log('tmp":', tmp);
          _.each(tmp.split(','), function(idString) {
            var id = parseInt(idString);
            var machine = $scope.machinesById[id];
            if (machine) {
              machineSkills.push(machine);
            }
          });
        }
        console.log('setting tutor.MachineSkills to ', machineSkills);
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

        var timeCreated = moment(p.Created);
        var timeStart = moment(p.TimeStart);
        var timeEnd = moment(p.TimeEnd);

        p.Created = timeCreated.tz('Europe/Berlin').format('YYYY-MM-DD');
        p.TimeStart = timeStart.tz('Europe/Berlin').format('YYYY-MM-DD HH:mm');
        p.TimeEnd = timeEnd.tz('Europe/Berlin').format('YYYY-MM-DD HH:mm');

        if (timeEnd.unix() > 0) {
          p.TimeEndLocal = timeEnd.tz('Europe/Berlin').format('YYYY-MM-DD HH:mm');
          var duration = timeEnd.clone().subtract(timeStart);
          p.TimerTimeTotalHours = duration.hours();
          p.TimerTimeTotalMinutes = duration.minutes();
        }

        if (timeEnd.unix() > 0) {
          var durationPlanned = timeEnd.clone().subtract(timeStart);
          p.ReservedTimeTotalHours = durationPlanned.hours();
          p.ReservedTimeTotalMinutes = durationPlanned.minutes();
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
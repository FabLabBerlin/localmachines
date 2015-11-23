(function(){

'use strict';

var app = angular.module('fabsmith.admin.tutoring', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/tutoring', {
    templateUrl: 'ng-modules/tutoring/tutoring.html',
    controller: 'TutoringCtrl'
  });
}]); // app.config

app.controller('TutoringCtrl', ['$scope', '$http', '$location', 
  function($scope, $http, $location) {

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
      console.log(settings);
      _.each(settings, function(setting) {
        $scope.settings[setting.Name] = setting;
      });
    })
    .error(function() {
      toastr.error('Failed to get global config');
    });
  };

  $scope.loadTutors = function() {
    $http({
      method: 'GET',
      url: '/api/tutoring/tutors',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(tutorList) {
      $scope.tutors = tutorList.Data;

      // TODO: Remember to remove this fake data generator
      _.each($scope.tutors, function(tutor) {
        tutor.Skills = "Laser Cutter, CNC Mill, MakerBot Replicatior";
      });
      console.log(tutorList);
    })
    .error(function() {
      toastr.error('Failed to load tutor list');
    });
  };

  $scope.loadPurchases = function() {
    $http({
      method: 'GET',
      url: '/api/tutoring/purchases',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(purchaseList) {
      $scope.purchases = purchaseList.Data;
      _.each($scope.purchases, function(purchase) {
        purchase.TutorName = 'Tina Atari';
        purchase.UserName = 'Milda Sane';

        purchase.Created = moment(purchase.StartTime).format('D MMM YY');
        purchase.TimeStart = moment(purchase.StartTime).format('D MMM YY HH:mm');
        purchase.TimeEnd = moment(purchase.EndTime).format('D MMM YY HH:mm');

        purchase.ReservedTimeTotalHours = purchase.TotalTime.toFixed(0);
        purchase.ReservedTimeTotalMinutes = 12;

        purchase.TimerTimeTotalHours = 0;
        purchase.TimerTimeTotalMinutes = 0;

        purchase.TimeTotal = purchase.TotalTime.toFixed(2);
      });
      console.log(purchaseList);
    })
    .error(function() {
      toastr.error('Failed to load purchase list');
    });
  };

  $scope.addTutor = function() {
    $location.path('/tutoring/tutor');
  };

  $scope.editTutor = function() {
    $location.path('/tutoring/tutor');
  };

  $scope.addPurchase = function() {
    $location.path('/tutoring/purchase');
  };

  $scope.editPurchase = function() {
    $location.path('/tutoring/purchase');
  };

  $scope.loadSettings();
  $scope.loadTutors();
  $scope.loadPurchases();

}]); // app.controller

})(); // closure
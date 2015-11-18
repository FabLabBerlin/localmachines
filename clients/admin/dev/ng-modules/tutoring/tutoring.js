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
      console.log(purchaseList);
    })
    .error(function() {
      toastr.error('Failed to load purchase list');
    });
  };

  $scope.addTutor = function() {
    $location.path('/tutoring/tutor');
  };

  $scope.addPurchase = function() {
    alert('Add purchase');
  };

  $scope.loadSettings();
  $scope.loadTutors();
  $scope.loadPurchases();

}]); // app.controller

})(); // closure
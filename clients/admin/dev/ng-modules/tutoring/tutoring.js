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

  $scope.tutors = [
    {Id:1, Name:'Ahmad Taleb', Price:60, PriceUnit:'hour'},
    {Id:2, Name:'Zelma Atari', Price:60, PriceUnit:'hour'}
  ];

  $scope.purchases = [
    {Id:1, TutorName:'Ahmad Taleb', UserName:'Selma Atari', StartTime:'15 Oct 2015 12:00', EndTime:'15 Oct 2015 14:00', TotalTime:1, TotalPrice:120, VAT:19},
    {Id:1, TutorName:'Ahmad Taleb', UserName:'Selma Atari', StartTime:'15 Oct 2015 16:00', EndTime:'15 Oct 2015 18:00', TotalTime:1, TotalPrice:120, VAT:19}
  ];

  // Calculate the TotalPriceVAT
  _.each($scope.purchases, function(purchase) {
    purchase.TotalPriceVAT = purchase.VAT * purchase.TotalPrice / 100.0;
  });

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
    alert('Add tutor');
  };

  $scope.addPurchase = function() {
    alert('Add purchase');
  };

  $scope.loadSettings();
  $scope.loadTutors();
  $scope.loadPurchases();

}]); // app.controller

})(); // closure
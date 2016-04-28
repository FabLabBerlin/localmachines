(function(){

'use strict';

var app = angular.module('fabsmith.admin.settings', 
 ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/settings', {
    templateUrl: 'ng-modules/settings/settings.html',
    controller: 'SettingsCtrl'
  });
}]); // app.config

app.controller('SettingsCtrl',
 ['$scope', '$http', '$location', '$cookies', 'randomToken',
 function($scope, $http, $location, $cookies, randomToken) {

  $scope.loadSettings = function() {
    $http({
      method: 'GET',
      url: '/api/settings',
      params: {
        location: $cookies.get('locationId'),
        ac: new Date().getTime()
      }
    })
    .success(function(settings) {
      $scope.settings = {
        Currency: {},
        TermsUrl: {},
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

  $scope.loadSettings();

  $scope.save = function() {
    $http({
      method: 'POST',
      url: '/api/settings?location=' + $cookies.get('locationId'),
      headers: {'Content-Type': 'application/json' },
      data: _.map($scope.settings, function(setting, name) {
        return _.extend({
          Name: name,
        }, setting);
      }),
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

}]); // app.controller

})(); // closure

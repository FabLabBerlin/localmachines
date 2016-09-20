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
 ['$scope', '$http', '$location', '$cookies', 'randomToken', 'api',
 function($scope, $http, $location, $cookies, randomToken, api) {

  api.loadSettings(function(settings) {
    $scope.settings = settings;
  });

  $scope.save = function() {
    $http({
      method: 'POST',
      url: '/api/settings?location=' + $cookies.get('location'),
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

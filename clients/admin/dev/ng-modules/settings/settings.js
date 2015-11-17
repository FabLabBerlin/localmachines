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

app.controller('SettingsCtrl', ['$scope', '$http', '$location', 'randomToken', 
 function($scope, $http, $location, randomToken) {

  $scope.loadSettings = function() {
    $http({
      method: 'GET',
      url: '/api/settings',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(settings){
      $scope.settings = settings;
    })
    .error(function(){
      toastr.error('Failed to get global config');
    });
  };

}]); // app.controller

})(); // closure

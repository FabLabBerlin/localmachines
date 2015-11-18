(function(){

'use strict';

var app = angular.module('fabsmith.admin.tutor', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/tutor', {
    templateUrl: 'ng-modules/tutor/tutor.html',
    controller: 'TutorCtrl'
  });
}]); // app.config

app.controller('TutorCtrl', ['$scope', '$http', '$location', 
  function($scope, $http, $location) {

  

}]); // app.controller

})(); // closure
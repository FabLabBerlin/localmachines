(function(){

'use strict';
var app = angular.module('fabsmith.admin.coworking', ['ngRoute', 'ngCookies', 'fabsmith.admin.randomtoken']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/coworking', {
    templateUrl: 'ng-modules/coworking/coworking.html',
    controller: 'CoworkingCtrl'
  });
}]); // app.config

app.controller('CoworkingCtrl',
 ['$scope', '$routeParams', '$http', '$location', 'randomToken',
 function($scope, $routeParams, $http, $location, randomToken) {

}]); // app.controller

})(); // closure

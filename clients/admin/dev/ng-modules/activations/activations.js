(function(){

'use strict';

var app = angular.module('fabsmith.admin.activations', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
    $routeProvider.when('/activations', {
        templateUrl: 'ng-modules/activations/activations.html',
        controller: 'ActivationsCtrl'
    });
}]); // app.config

app.controller('ActivationsCtrl', ['$scope', '$http', '$location', 
 function($scope, $http, $location) {

    $('.datepicker').pickadate();
    $('.timepicker').pickatime();

}]); // app.controller


})(); // closure
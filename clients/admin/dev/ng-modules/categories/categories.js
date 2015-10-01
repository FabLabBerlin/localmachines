(function(){

'use strict';

var app = angular.module('fabsmith.admin.categories', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/categories', {
    templateUrl: 'ng-modules/categories/categories.html',
    controller: 'CategoriesCtrl'
  });
}]); // app.config

app.controller('CategoriesCtrl', ['$scope', '$http', '$location', '$cookieStore', 
 function($scope, $http, $location, $cookieStore) {


}]); // app.controller

})();
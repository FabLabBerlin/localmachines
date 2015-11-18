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
    {Id:1, TutorName:'Ahmad Taleb', UserName:'Selma Atari', StartTime:'15 Oct 2015 12:00', EndTime:'15 Oct 2015 14:00', TotalTime:1, TotalPrice:120, VAT:22.8},
    {Id:1, TutorName:'Ahmad Taleb', UserName:'Selma Atari', StartTime:'15 Oct 2015 16:00', EndTime:'15 Oct 2015 18:00', TotalTime:1, TotalPrice:120, VAT:22.8}
  ];

  $scope.addTutor = function() {
    alert('Add tutor');
  }

  $scope.addPurchase = function() {
    alert('Add purchase');
  }

}]); // app.controller

})(); // closure
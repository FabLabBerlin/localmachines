(function(){

'use strict';

var app = angular.module('fabsmith.admin.tutoring.tutor', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/tutoring/tutor', {
    templateUrl: 'ng-modules/tutoring/tutor.html',
    controller: 'TutorCtrl'
  });
}]); // app.config

app.controller('TutorCtrl', ['$scope', '$http', '$location', 
  function($scope, $http, $location) {

  $scope.tutor = {
    PriceUnit: 'hour'
  };

  $scope.cancel = function() {
    $location.path('/tutoring');
  };

  $scope.save = function() {

    if ($scope.tutor.Name === '' || !$scope.tutor.Name) {
      alert('Enter tutor name');
      return;
    }

    if ($scope.tutor.Price === '' || !$scope.tutor.Price) {
      alert('Enter tutor price');
      return;
    }

    console.log($scope.tutor);

    var tutor = $scope.tutor;
    tutor.Price = parseInt(tutor.Price);

    $http({
      method: 'PUT',
      url: '/api/tutoring/tutor',
      data: tutor,
      headers: { 
        'Content-Type': 'application/json' 
      },
      params: { 
        ac: new Date().getTime() 
      },
      transformRequest: function(data) {
        return JSON.stringify(data);
      }
    })
    .success(function(data) {
      toastr.success('Tutor updated');
      $location.path('/tutoring');
    })
    .error(function() {
      toastr.error('Failed to save tutor');
    });
  };

}]); // app.controller

})(); // closure
(function(){

'use strict';

var app = angular.module("fabsmith.admin.reservations.toggle", []);

app.directive('toggle', function() {
  return {
    restrict: 'E',
    templateUrl: 'ng-modules/reservations/toggle.html',
    transclude: true,
    scope: true,
    controller: function($scope, $element, $attrs, $http, $location) {
      $scope.title = $attrs.title;
      $scope.show = false;
      $scope.toggle = function() {
        console.log('toogle', $scope.show);
        $scope.show = !$scope.show;
      };
    },
    link: function($scope, $element, $attrs) {
    }
  };
});

})(); // closure

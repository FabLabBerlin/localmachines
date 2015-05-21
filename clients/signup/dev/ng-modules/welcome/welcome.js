(function(){

'use strict';

angular.module('fabsmith.signup.welcome', ['ngRoute'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/welcome', {
    templateUrl: 'ng-modules/welcome/welcome.html',
    controller: 'WelcomeCtrl'
  });
}])

.controller('WelcomeCtrl', ['$scope', '$location',
 function($scope, $location) {

   $scope.goToSignUp = function() {
     $location.path('/form');
   };

}]);

})();

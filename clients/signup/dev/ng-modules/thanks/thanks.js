(function(){

'use strict';

angular.module('fabsmith.signup.thanks', ['ngRoute'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/thanks', {
    templateUrl: 'ng-modules/thanks/thanks.html',
    controller: 'ThanksCtrl'
  });
}])

.controller('ThanksCtrl', ['$scope', '$location', 
 function($scope, $location) {
  
  $scope.backToForm = function() {
    $location.path('/form');
  };

}]);

})();
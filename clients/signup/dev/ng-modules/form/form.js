(function(){

'use strict';

angular.module('fabsmith.signup.form', ['ngRoute'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/form', {
    templateUrl: 'ng-modules/form/form.html',
    controller: 'FormCtrl'
  });
}])

.controller('FormCtrl', ['$scope', '$location', 
 function($scope, $location) {
  
  $scope.submitForm = function() {
    $location.path('/thanks');
  };

}]);

})();
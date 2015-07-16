(function(){

'use strict';

angular.module('fabsmith.signup.payment', ['ngRoute'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/payment', {
    templateUrl: 'ng-modules/payment/payment.html',
    controller: 'PaymentCtrl'
  });
}])

.controller('PaymentCtrl', ['$scope', '$location', '$http',
  function($scope, $location, $http) {

    $.loadPublicKey = function() {
      $http({
        method: 'GET',
        url: '/api/paymill/public'
      })
      .success(function(data) {
        $scope.public_key = data.PublicKey;
      })
      .error(function() {
        toastr.error('Error while trying to register');
      });
    };

    console.log("Payment!");
    $.loadPublicKey();

  }]);
})();

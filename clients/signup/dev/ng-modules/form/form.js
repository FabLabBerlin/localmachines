(function(){

'use strict';

angular.module('fabsmith.signup.form', ['ngRoute'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/form', {
    templateUrl: 'ng-modules/form/form.html',
    controller: 'FormCtrl'
  });
}])

.controller('FormCtrl', ['$scope', '$location', '$http',
 function($scope, $location, $http) {

  $scope.emailRegExp = /^[_a-z0-9]+(\.[_a-z0-9]+)*@[a-z0-9-]+(\.[a-z0-9-]+)*(\.[a-z]{2,4})$/;
  $scope.minUsernameAndPasswordLength = 3;

  $scope.submitForm = function() {

    // Check if email is empty
    if(!$scope.email || $scope.email === ''){
      toastr.error('Please enter an e-mail');
      return;
    }

    // Check if E-Mail is valid
    if(!$scope.emailRegExp.test($scope.email)){
      toastr.error('Please enter a valid e-mail');
      return;
    }

    // Check if First Name is empty
    if(!$scope.firstName || $scope.firstName === ''){
      toastr.error('Please enter a first name');
      return;
    }

    // Check if Last Name is empty
    if(!$scope.lastName || $scope.lastName === ''){
      toastr.error('Please enter a last name');
      return;
    }

    // Check if Username is empty
    if(!$scope.username || $scope.username === ''){
      toastr.error('Please enter an username');
      return;
    }

    // Check if Username is long enough
    if($scope.username.length < $scope.minUsernameAndPasswordLength){
      toastr.error('Please enter an username which is at least 3 characters long');
      return;
    }

    // Check if Password is empty
    if(!$scope.password || $scope.password === ''){
      toastr.error('Please enter a password');
      return;
    }

    // Check if Password is long enough
    if($scope.password.length < $scope.minUsernameAndPasswordLength){
      toastr.error('Please enter a password which is at least 3 characters long');
      return;
    }

    // Check if password matches the repeated password
    if ($scope.password !== $scope.passwordRepeat) {
      toastr.error('Passwords do not match');
      return;
    }


    $http({
      method: 'POST',
      url: '/api/users/signup',
      data: {
        "User": {
          "Email": $scope.email,
          "Company": $scope.company,
          "FirstName": $scope.firstName,
          "LastName": $scope.lastName,
          "Username": $scope.username
        },
        "Password": $scope.password
      },
      transformRequest: function(data) {
        return JSON.stringify(data);
      },
      headers: {'Content-Type': 'application/json'}
    })
    .success(function() {
      toastr.success("Registration successful !");
      $location.path('/thanks');
    })
    .error(function() {
      toastr.error('Error while trying to register');
    });

  };

}]);

})();

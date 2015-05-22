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

    // Check if AGB and Data Protection Agreement is checked
    if(!$scope.agb_dpa_agreed){
      toastr.error('You have to agree to the AGB and Data Protection Agreement');
      return;
    }

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
        "Password": $scope.password,
        "RegistrationDate" : Date.now()
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

  var numberOfSecondsBeforeResting = 5;
  // Instant call function
  (function(nbSecToReset){
    // Check for idle time
    var idleTime = 0;
    $(document).ready(function () {
      //Increment the idle time counter every minute.
      var idleInterval = setInterval(timerIncrement, /*60*/ 1 * 1000); // 1 minute

      //Zero the idle timer on mouse movement or a keypress.
      $(this).mousemove(function (e) {
        idleTime = 0;
      });
      $(this).keypress(function (e) {
        idleTime = 0;
      });
    });

    function timerIncrement() {
      idleTime = idleTime + 1;
      // After 5 minutes we reset the form
      if (idleTime > (nbSecToReset-1)) {
        $scope.email = "";
        $scope.company = "";
        $scope.firstName = "";
        $scope.lastName = "";
        $scope.username = "";
        $scope.password = "";
        $scope.$apply();
        idleTime = 0;
      }
    }
  })(numberOfSecondsBeforeResting);

  $scope.displayAGB = function(){
    vex.defaultOptions.className = 'vex-theme-wireframe';
    vex.dialog.alert(
      // "<p>Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.</p>" +
      // "<p>Duis autem vel eum iriure dolor in hendrerit in vulputate velit esse molestie consequat, vel illum dolore eu feugiat nulla facilisis at vero eros et accumsan et iusto odio dignissim qui blandit praesent luptatum zzril delenit augue duis dolore te feugait nulla facilisi. Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat.</p>" +
      // "<p>Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat. Duis autem vel eum iriure dolor in hendrerit in vulputate velit esse molestie consequat, vel illum dolore eu feugiat nulla facilisis at vero eros et accumsan et iusto odio dignissim qui blandit praesent luptatum zzril delenit augue duis dolore te feugait nulla facilisi.</p>" +
      "<p>Nam liber tempor cum soluta nobis eleifend option congue nihil imperdiet doming id quod mazim placerat facer possim assum. Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat. Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.</p>"
    );
  };

  $scope.displayDatenschutzbestimmungen = function(){
    vex.defaultOptions.className = 'vex-theme-wireframe';
    vex.dialog.alert(
      "<p>Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.</p>" +
      "<p>Duis autem vel eum iriure dolor in hendrerit in vulputate velit esse molestie consequat, vel illum dolore eu feugiat nulla facilisis at vero eros et accumsan et iusto odio dignissim qui blandit praesent luptatum zzril delenit augue duis dolore te feugait nulla facilisi. Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat.</p>" +
      "<p>Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat. Duis autem vel eum iriure dolor in hendrerit in vulputate velit esse molestie consequat, vel illum dolore eu feugiat nulla facilisis at vero eros et accumsan et iusto odio dignissim qui blandit praesent luptatum zzril delenit augue duis dolore te feugait nulla facilisi.</p>" +
      "<p>Nam liber tempor cum soluta nobis eleifend option congue nihil imperdiet doming id quod mazim placerat facer possim assum. Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat. Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.</p>"
    );
  };

}]);
})();

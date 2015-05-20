(function(){

'use strict';

// Declare app level module which depends on views, and components
var app = angular.module('fabsmith.signup', [
  'ngRoute',
  'fabsmith.signup.form',
  'fabsmith.signup.thanks'
]);

// Configure http provider to send regular form POST data instead of JSON
app.config(['$httpProvider', function($httpProvider) {
  $httpProvider.defaults.transformRequest = function(data){
    if (data === undefined) {
      return data;
    }
    return $.param(data);
  };
  $httpProvider.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded; charset=UTF-8';
}]);

// Signup main controller
app.controller('MainCtrl', ['$scope', '$location', 
 function($scope, $location){

  // Configure toastr default location
  toastr.options.positionClass = 'toast-bottom-left';
  toastr.info('MainCtrl loaded');

  // Configure vex theme
  vex.defaultOptions.className = 'vex-theme-custom';

  // Redirect
  $location.path('/form');

}]);

})(); // closure

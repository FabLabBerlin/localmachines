(function(){

'use strict';

var app = angular.module('fabsmith.admin.mainmenu', ['ngRoute', 'ngCookies']);

app.directive('mainMenu', function() {
  return {
    templateUrl: 'ng-modules/mainmenu/mainmenu.html',
    restrict: 'E',
    controller: [
      '$scope', '$element', '$cookieStore', 
      function($scope, $element, $cookieStore) {

      $scope.userFullName = $cookieStore.get('FirstName') +
    ' ' + $cookieStore.get('LastName');

      var links = $($element).find('a');
      links.click(function(){

        // The float CSS parameter is changed to 'none' whenever the
        // window width is below Bootstrap grid breakpoint (768px)
        var navfloat = $('.navbar-header').css('float');
        if (navfloat === 'none') {
          $($element).find('.navbar-collapse').collapse('hide');
        }       
      });
    }]
  };
});

})(); // closure
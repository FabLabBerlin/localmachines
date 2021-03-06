(function(){

'use strict';

var app = angular.module('fabsmith.admin.mainmenu', ['ngRoute', 'ngCookies']);

app.directive('mainMenu', function() {
  return {
    templateUrl: 'ng-modules/mainmenu/mainmenu.html',
    restrict: 'E',
    controller: ['$rootScope', '$scope', '$element', '$cookies', '$http',
                 function($rootScope, $scope, $element, $cookies, $http) {
      $scope.data = $rootScope.mainMenu;
      $scope.location = $rootScope.location;
      var links = $($element).find('a');
      links.click(function(){
        if ($(this).attr('data-toggle') === 'dropdown') {
          return;
        }
        // The float CSS parameter is changed to 'none' whenever the
        // window width is below Bootstrap grid breakpoint (768px)
        var navfloat = $('.navbar-header').css('float');
        if (navfloat === 'none') {
          $($element).find('.navbar-collapse').collapse('hide');
        }
      });
    }],
    link: function($rootScope, $scope, $element, $attrs) {}
  };
});

})(); // closure
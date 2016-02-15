(function(){

'use strict';

var app = angular.module('fabsmith.admin.mainmenu', ['ngRoute', 'ngCookies']);

app.directive('mainMenu', function() {
  return {
    templateUrl: 'ng-modules/mainmenu/mainmenu.html',
    restrict: 'E',
    controller: ['$rootScope', '$scope', '$element', '$cookieStore', function($rootScope, $scope, $element, $cookieStore) {
      $scope.data = $rootScope.mainMenu;

      var links = $($element).find('a');
      links.click(function(){

        // The float CSS parameter is changed to 'none' whenever the
        // window width is below Bootstrap grid breakpoint (768px)
        var navfloat = $('.navbar-header').css('float');
        if (navfloat === 'none') {
          $($element).find('.navbar-collapse').collapse('hide');
        }
      });
    }],
    link: function($rootScope, $scope, $element, $attrs) {
      console.log('main menu attrs: ', $attrs);
    }
  };
});

})(); // closure
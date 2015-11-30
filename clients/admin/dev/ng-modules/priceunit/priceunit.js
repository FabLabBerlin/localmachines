(function(){

'use strict';

angular.module("fabsmith.admin.priceunit", [])
.directive('priceunit', function() {
  return {
    restrict: 'E',
    scope: { product: '=' },
    templateUrl: 'ng-modules/priceunit/priceunit.html'
  };
});

})(); // closure

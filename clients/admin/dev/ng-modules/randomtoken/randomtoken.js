var mod = angular.module("fabsmith.admin.randomtoken", []);

mod.directive('randomToken', function() {
  return {
    templateUrl: 'ng-modules/randomtoken/randomtoken.html',
    restrict: 'E',
    controller: ['$scope', '$element', function($scope, $element){
      $scope.randomToken = 'RandomToksmen';

      $scope.generateRandomToken = function() {
        var tokens = [
          'Randy3time',
          'Token2be4me',
          'Token4life',
          'Randomi7er',
          'RandomSk8ter',
          'H8tersGonn4'
        ];

        var id = Math.round(Math.random() * tokens.length);

        $scope.randomToken = tokens[id];
      };
    }]
  };
});
(function(){

'use strict';

var app = angular.module('fabsmith.admin.mainmenu', ['ngRoute', 'ngCookies']);

app.controller('MainMenuCtrl', ['$scope', '$http', '$location', function($scope, $http, $location) {
	//
}]); // app.controller

app.directive('mainMenu', function() {
	return {
		templateUrl: 'ng-modules/mainmenu/mainmenu.html',
		restrict: 'E',
		controller: ['$scope', function($scope){
			//
		}]
	};
});

})(); // closure
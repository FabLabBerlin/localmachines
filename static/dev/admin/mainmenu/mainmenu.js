(function(){

'use strict';

var app = angular.module('fabsmith.backoffice.mainmenu', ['ngRoute', 'ngCookies']);

app.controller('MainMenuCtrl', ['$scope', '$http', '$location', function($scope, $http, $location) {
	//
}]); // app.controller

app.directive('mainMenu', function() {
	return {
		templateUrl: 'static/dev/admin/mainmenu/mainmenu.html',
		restrict: 'E',
		controller: ['$scope', function($scope){
			//
		}]
	}
});

})(); // closure
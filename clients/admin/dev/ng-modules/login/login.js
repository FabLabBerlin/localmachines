(function(){

'use strict';

var app = angular.module('fabsmith.admin.login', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/login', {
    templateUrl: 'ng-modules/login/login.html',
    controller: 'LoginCtrl'
  });
}]); // app.config

app.controller('LoginCtrl', ['$scope', '$http', '$location', function($scope, $http, $location) {
  // Local login function - if we do it by entering username and password in the browser
  if (window.libnfc) {
    $scope.nfcSupport = true;

    $scope.onNfcError = function(error) {
      window.libnfc.cardRead.disconnect($scope.loginWithUid);
      window.libnfc.cardReaderError.disconnect($scope.onNfcError);
      toastr.error(error);
      $scope.nfcErrorTimeout = setTimeout($scope.getNfcUid, 2000);
    };

    $scope.loginWithUid = function(uid) {
      window.libnfc.cardRead.disconnect($scope.loginWithUid);
      window.libnfc.cardReaderError.disconnect($scope.onNfcError);
      $http({
        method: 'POST',
        url: '/api/users/loginuid',
        data: {
          uid: uid,
          ac: new Date().getTime()
        }
      })
      .success(function(data) {
        if (data.UserId) {
          $scope.getUserData(data.UserId);
        }
      })
      .error(function() {
        toastr.error('Failed to log in');
        setTimeout($scope.getNfcUid, 2000);
      });
    };

    $scope.getNfcUid = function() {
      //window.libnfc.nfcReaderError.connect($scope.onNfcError);
      window.libnfc.cardRead.connect($scope.loginWithUid);
      window.libnfc.cardReaderError.connect($scope.onNfcError);
      window.libnfc.asyncScan(); // For infinite amount of time
    };

    $scope.getNfcUid();
  } // if window.libnfc

  $scope.getUserData = function(userId) {
    $http({
      method: 'GET',
      url: '/api/users/' + userId,
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data){
      $scope.$emit('user-login', data);
      $location.path('/dashboard');
    })
    .error(function(data, status){
      console.log('Status: ' + status);
      console.log('Data' + data);
      toastr.error('Could not get user data');
    });
  };

  $scope.login = function() {
    if (window.libnfc) {
      try {
        window.libnfc.cardRead.disconnect($scope.loginWithUid);
      } catch (err) {
        toastr.error(err.message);
      }

      try {
        window.libnfc.cardReaderError.disconnect($scope.onNfcError);
      } catch (err) {
        toastr.error(err.message);
      }

      clearTimeout($scope.nfcErrorTimeout);
    }

    $http({
      method: 'POST',
      url: '/api/users/login',
      data: {
        username: $scope.username,
        password: $scope.password,
      },
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      if (data.UserId) {
        $scope.getUserData(data.UserId);
      }
    })
    .error(function() {
      toastr.error('Failed to log in');
    });
  };
  
}]); // app.controller

})(); // closure
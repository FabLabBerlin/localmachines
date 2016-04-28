(function(){

'use strict';

var app = angular.module('fabsmith.admin.login', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/login', {
    templateUrl: 'ng-modules/login/login.html',
    controller: 'LoginCtrl'
  });
}]); // app.config

app.controller('LoginCtrl',
 ['$rootScope', '$scope', '$http', '$location', '$cookies',
 function($rootScope, $scope, $http, $location, $cookies) {
  // Local login function - if we do it by entering username and password in the browser
  if (window.libnfc) {
    $scope.nfcSupport = true;
    $scope.locations = [];

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
          $cookies.put('locationId', 1);
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
      //$scope.$emit('user-login', data);
      $rootScope.mainMenu.userFullName = data.FirstName + ' ' + data.LastName;
      $location.path('/machines');
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

    var locationId = $('select[name="location"]').val();
    locationId = parseInt(locationId);

    $http({
      method: 'POST',
      url: '/api/users/login',
      data: {
        username: $scope.username,
        password: $scope.password,
        location: locationId,
        admin: true
      },
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      if (data.UserId) {
        $cookies.put('locationId', locationId);
        $scope.getUserData(data.UserId);
      }
    })
    .error(function() {
      toastr.error('Failed to log in');
    });
  };

  $scope.getLocations = function() {
    $http({
      method: 'GET',
      url: '/api/locations',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(locations) {
      $scope.locations = locations;
    })
    .error(function() {
      toastr.error('Failed to load locations');
    });
  };

  $scope.getLocations();
  
}]); // app.controller

})(); // closure
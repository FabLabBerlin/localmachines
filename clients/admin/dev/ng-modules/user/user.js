(function(){

'use strict';
var app = angular.module('fabsmith.admin.user', ['ngRoute', 'ngCookies', 'fabsmith.admin.randomtoken']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/users/:userId', {
    templateUrl: 'ng-modules/user/user.html',
    controller: 'UserCtrl'
  });
}]); // app.config

app.controller('UserCtrl', 
 ['$scope', '$routeParams', '$http', '$location', 'randomToken', 
 function($scope, $routeParams, $http, $location, randomToken) {

  // Check for NFC browser
  if (window.libnfc) {
    $scope.nfcSupport = true;
    $scope.nfcPolling = false;
    $scope.nfcButtonLabel = "Read NFC UID";

    $scope.resetNfcUi = function() {
      $scope.nfcPolling = false;
      $scope.nfcButtonLabel = "Read NFC UID";
      $scope.$apply();
    };

    $scope.onNfcUid = function(uid) {
      window.libnfc.cardRead.disconnect($scope.onNfcUid);
      window.libnfc.cardReaderError.disconnect($scope.onNfcError);
      clearTimeout($scope.getNfcUidTimeout);
      $scope.nfcUid = uid;
      $scope.resetNfcUi();
    };    

    // Cancel callback
    $scope.cancelGetNfcUid = function() {
      window.libnfc.cardRead.disconnect($scope.onNfcUid);
      window.libnfc.cardReaderError.disconnect($scope.onNfcError);
      toastr.warning("Reading NFC took too long");
      $scope.resetNfcUi();
    };

    $scope.onNfcError = function(error) {
      window.libnfc.cardRead.disconnect($scope.onNfcUid);
      window.libnfc.cardReaderError.disconnect($scope.onNfcError);
      clearTimeout($scope.getNfcUidTimeout);
      toastr.error(error);
      $scope.resetNfcUi();
    };

    $scope.getNfcUid = function() {
      // Add event listener
      window.libnfc.cardRead.connect($scope.onNfcUid);
      window.libnfc.cardReaderError.connect($scope.onNfcError);

      // Start waiting for the NFC card to approach the reader
      window.libnfc.asyncScan();
      $scope.nfcPolling = true;
      $scope.nfcButtonLabel = "Waiting for NFC card...";

      // Cancel scan after timeout
      $scope.getNfcUidTimeout = setTimeout($scope.cancelGetNfcUid, 10000);
    };
  }
  
  // Init scope variables
  var pickadateOptions = {
    format: 'yyyy-mm-dd'
  };
  $('.datepicker').pickadate(pickadateOptions);
  $scope.user = { Id: $routeParams.userId };
  $scope.userMachines = [];
  $scope.userMemberships = [];

  $scope.loadUserData = function() {
    $http({
      method: 'GET',
      url: '/api/users/' + $scope.user.Id,
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(user) {
      $scope.user = user;
      if (user.UserRole === 'admin') {
        user.Admin = true;
      }
      $scope.loadAvailableMachines();
    })
    .error(function(data, status) {
      toastr.error('Failed to load user data');
    });
  };

  // It should go like this:
  //
  // 1. Load user data
  // 2. Load available machines
  // 3. Load user machine permissions
  // 4. Load available memberships
  // 5. Load user memberships

  $scope.loadUserData();

  $scope.loadAvailableMachines = function() {
    $http({
      method: 'GET',
      url: '/api/machines',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(availableMachines) {
      $scope.availableMachines = availableMachines;

      if ($scope.user.UserRole === 'admin') {
        _.each($scope.availableMachines, function(machine){
          machine.Disabled = true;
          machine.Checked = true;
        });
        $scope.getAvailableMemberships();
      } else {
        $scope.loadUserMachinePermissions($scope.getAvailableMemberships);
      }
      
    })
    .error(function() {
      console.log('Could not get machines');
    });
  };

  $scope.loadUserMachinePermissions = function(callback) {
    $http({
      method: 'GET',
      url: '/api/users/' + $scope.user.Id + '/machines',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(userMachines) {

      $scope.userMachines = userMachines;

      if (callback) {
        callback();
      }
      
      _.each($scope.availableMachines, function(machine) {
        machine.Checked = false;
        _.each(userMachines, function(userMachine) {
          if (userMachine.Id === machine.Id) {
            machine.Checked = true;
          }
        }); // each
      }); // each
    }) // success
    .error(function(msg, status) {
      console.log(msg);
      console.log('Could not get user machines');
    });
  };

  // TODO: this could be transformed into a filter
  function formatDate(d) {
    //console.log('formatDate: d: ', d);
    var mm = (d.getMonth() + 1);
    if (mm < 10) {
      mm = '0' + mm;
    }
    var dd = d.getDate();
    if (dd < 10) {
      dd = '0' + dd;
    }
    return d.getFullYear() + '-' + mm + '-' + dd;
  }

  $scope.getAvailableMemberships = function() {
    $http({
      method: 'GET',
      url: '/api/memberships',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      $scope.memberships = data;
      $scope.membershipsById = {};
      _.each($scope.memberships, function(m) {
        $scope.membershipsById[m.Id] = m;
      });
      
      $scope.getUserMemberships();
    })
    .error(function(data, status) {
      console.log('Could not get memberships');
      console.log('Data: ' + data);
      console.log('Status code: ' + status);
    });
  };

  $scope.getUserMemberships = function() {
    $http({
      method: 'GET',
      url: '/api/users/' + $scope.user.Id + '/memberships',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      $scope.userMemberships = _.map(data, function(userMembership) {
        userMembership.StartDate = new Date(Date.parse(userMembership.StartDate));
        _.merge(userMembership, {
          EndDate: new Date(userMembership.StartDate)
        });
        //console.log('userMembership.Id: ', userMembership.Id);
        var membership = $scope.membershipsById[userMembership.MembershipId];
        //console.log('membership: ', membership);
        userMembership.EndDate.setDate(userMembership.StartDate.getDate() + 
         membership.Duration); 
        userMembership.StartDate = formatDate(userMembership.StartDate);
        userMembership.EndDate = formatDate(userMembership.EndDate);
        return userMembership;
      });
    })
    .error(function(data, status) {
      console.log('Could not get user memberships');
      console.log('Data: ' + data);
      console.log('Status code: ' + status);
    });
  };

  $scope.addUserMembership = function() {
    var startDate = $('#adm-add-user-membership-start-date').val();
    if (!startDate) {
      toastr.error('Please select a Start Date');
      return;
    }

    if ($scope.overlapsUserMembership(startDate)) {
      toastr.error('Overlapping existing membership');
      return;
    }

    var userMembershipId = $('#user-select-membership').val();
    if (!userMembershipId) {
      toastr.error('Please select a Membership');
      return;
    }
    $http({
      method: 'POST',
      url: '/api/users/' + $scope.user.Id + '/memberships',
      data: {
        StartDate: startDate,
        UserMembershipId: userMembershipId
      }, 
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function() {
      toastr.success('Membership created');
      $scope.getUserMemberships();
    })
    .error(function() {
      toastr.error('Error while trying to create new User Membership');
    });
  };

  $scope.overlapsUserMembership = function(startDateStr) {
    var startDate = new Date(startDateStr);
    var overlap = false;
    _.each($scope.userMemberships, function(mbs) {
      var endDate = new Date(mbs.EndDate);
      if (startDate <= endDate){
        overlap = true;
      }
    });
    return overlap;
  };

  $scope.cancel = function() {
    if (confirm('All changes will be discarded, click ok to continue.')) {
      $location.path('/users');
    }
  };

  $scope.deleteUserPrompt = function() {
    var token = randomToken.generate();
    vex.dialog.prompt({
      message: 'Enter <span class="delete-prompt-token">' + 
       token + '</span> to delete',
      placeholder: 'Token',
      callback: $scope.deleteUserPromptCallback.bind(this, token)
    });
  };

  $scope.deleteUserPromptCallback = function(expectedToken, value) {
    if (value) {    
      if (value === expectedToken) {
        $scope.deleteUser();
      } else {
        toastr.error('Wrong token');
      }
    } else if (value !== false) {
      toastr.error('No token');
    }
  };

  $scope.deleteUser = function() {
    $http({
      method: 'DELETE',
      url: '/api/users/' + $scope.user.Id,
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function() {
      toastr.success('User deleted');
      $location.path('/users');
    })
    .error(function() {
      toastr.error('Error while trying to delete user');
    });
  };

  $scope.deleteUserMembershipPrompt = function(userMembershipId) {
    var token = randomToken.generate();
    vex.dialog.prompt({
      message: 'Enter <span class="delete-prompt-token">' +
      token + '</span> to delete',
      placeholder: 'Token',
      callback: function(value) {
        if (value) {    
          if (value === token) {
            $scope.deleteUserMembership(userMembershipId);
          } else {
            toastr.error('Wrong token');
          }
        } else if (value !== false) {
          toastr.error('No token');
        }
      } // callback
    });
  };

  $scope.deleteUserMembership = function(userMembershipId) {
    $http({
      method: 'DELETE',
      url: '/api/users/' + $scope.user.Id + '/memberships/' + userMembershipId,
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      toastr.success('Membership deleted');
      //window.location.reload(true);
      $scope.getUserMemberships();
    })
    .error(function() {
      toastr.error('Error while trying to delete user membership');
    });
  };

  $scope.saveUser = function() {
    if ($scope.user.UserRole === 'admin') {
      $scope.updateUser();
    } else {
      $scope.updateUserMachinePermissions($scope.updateUser);
    }
  };

  $scope.updateUser = function(callback) {

    if ($scope.user.Admin) {
      $scope.user.UserRole = 'admin';
    } else {
      $scope.user.UserRole = '';
    }

    $http({
      method: 'PUT',
      url: '/api/users/' + $scope.user.Id,
      headers: {'Content-Type': 'application/json' },
      data: {
        User: $scope.user
      },
      transformRequest: function(data) {
        return JSON.stringify(data);
      },
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function() {
      if (callback) {
        callback();
      }
      toastr.success('User updated');
    })
    .error(function(data) {
      if (data === 'lastAdmin') {
        $scope.user.Admin = true;
        $scope.updateAdminStatus();
        toastr.error('You are the last remaining admin. Remember - power comes with great responsibility!');
      } else if (data === 'selfAdmin') { 
        $scope.user.Admin = true;
        $scope.updateAdminStatus();
        toastr.error('You can not unadmin yourself. Someone else has to do it.');
      } else {
        toastr.error('Error while trying to save changes');
      }
    });
  };

  $scope.updateUserMachinePermissions = function(callback) {
    var permissions = [];
    _.each($scope.availableMachines, function(machine) {
      if (machine.Checked) {
        var p = {MachineId: machine.Id};
        permissions.push(p);
      }
    });

    $http({
      method: 'PUT',
      url: '/api/users/' + $scope.user.Id + '/permissions',
      headers: {'Content-Type': 'application/json' },
      data: permissions,
      transformRequest: function(data) {
        return JSON.stringify(data);
      },
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function() {
      if (callback) {
        callback();
      }
    })
    .error(function() {
      toastr.error('Error while trying to update permissions');
    });
  };

  $scope.updatePassword = function() {
    
    // Check user entered password
    var minPassLength = 3;

    // If there is password at all
    if (!$scope.userPassword || $scope.userPassword === '') {
      toastr.error('No password');
      return;
    }

    // If it is long enough
    if ($scope.userPassword.length < minPassLength) {
      toastr.error('Password too short');
      return;
    }

    // If it matches the repeated password
    if ($scope.userPassword !== $scope.userPasswordRepeat) {
      toastr.error('Passwords do not match');
      return;
    }

    $http({
      method: 'POST', // TODO: change this to PUT as it is an update operation?
      url: '/api/users/' + $scope.user.Id + '/password',
      params: {
        password: $scope.userPassword,
        ac: new Date().getTime()
      }
    })
    .success(function() {
      toastr.success('Password successfully updated');
    })
    .error(function() {
      toastr.error('Error while trying to update password');
    });
  };

  $scope.updateNfcUid = function() {
    var minUidLen = 4;

    // If there is an uid at all
    if (!$scope.nfcUid || $scope.nfcUid === '') {
      toastr.error('No NFC UID');
      return;
    }

    // If it is long enough
    if ($scope.nfcUid.length < minUidLen) {
      toastr.error('NFC UID too short');
      return;
    }

    $http({
      method: 'PUT',
      url: '/api/users/' + $scope.user.Id + '/nfcuid',
      params: {
        nfcuid: $scope.nfcUid,
        ac: new Date().getTime()
      }
    })
    .success(function() {
      toastr.success('NFC UID successfully updated');
    })
    .error(function() {
      toastr.error('Error while updating NFC UID');
    });
  };

  $scope.updateUserRoles = function(value) {
    console.log('value:',value);
  };

  $scope.updateAdminStatus = function() {
    if ($scope.user.Admin) {
      _.each($scope.availableMachines, function(machine){
        machine.Checked = true;
        machine.Disabled = true;
      });
    } else {
      _.each($scope.availableMachines, function(machine){
        machine.Checked = false;
        machine.Disabled = false;
      });
      $scope.loadUserMachinePermissions();
    }
  };
}]); // app.controller

})(); // closure
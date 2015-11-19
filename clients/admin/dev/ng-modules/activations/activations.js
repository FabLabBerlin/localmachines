(function(){

'use strict';

var app = angular.module('fabsmith.admin.activations', 
 ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/activations', {
    templateUrl: 'ng-modules/activations/activations.html',
    controller: 'ActivationsCtrl'
  });
}]); // app.config

app.controller('ActivationsCtrl', ['$scope', '$http', '$location', 'randomToken', 
 function($scope, $http, $location, randomToken) {

  // We need full machine names for the activation table
  $scope.loadMachines = function() {
    $http({
      method: 'GET',
      url: '/api/machines',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(machines){
      $scope.machines = machines;
    })
    .error(function(){
      toastr.error('Failed to get machines');
    });
  };

  // Load machines before doing anything else
  if (!$scope.machines) {
    $scope.loadMachines();
  }

  // Loads and reloads activations according to filter.
  // If user ID is not set - load all users
  $scope.loadActivations = function() {
    
    if (!$scope.machines) {
      toastr.error('Machines not loaded');
      return;
    }

    $http({
      method: 'GET',
      url: '/api/activations', 
      params: {
        startDate: $scope.activationsStartDate,
        endDate: $scope.activationsEndDate,
        userId: 1,
        includeInvoiced: false,
        itemsPerPage: $scope.itemsPerPage,
        page: $scope.currentPage,
        ac: new Date().getTime()
      }
    })
    .success(function(response) {
      console.log('response:', response);
      _.each(response.ActivationsPage, function(activation){
        var machine = _.find($scope.machines, 'Id', activation.MachineId);
        if (machine) {
          activation.MachineName = machine.Name;
        } else {
          activation.MachineName = 'Undefined';
        }
      });

      var uniqUserIds = _.pluck(_.uniq(response.ActivationsPage, 'UserId'), 'UserId');
      _.each(uniqUserIds, $scope.loadUserName);

      $scope.activations = response.ActivationsPage;
      $scope.numActivations = response.NumActivations;
      $scope.numPages = Math.ceil($scope.numActivations / $scope.itemsPerPage);
    })
    .error(function() {
      toastr.error('Failed to load activations');
    });
  };

  $scope.loadUserName = function(userId) {
    $http({
      method: 'GET',
      url: '/api/users/' + userId + '/name',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(user) {
      _.each($scope.activations, function(activation) {
        if (activation.UserId === user.UserId) {
          activation.UserName = user.FirstName + ' ' + user.LastName;
        }
      });
    })
    .error(function() {
      toastr.error('Failed to load user name');
    });
  };

  // This is called whenever start or end date changes
  $scope.onFilterChange = function() {

    // Check if both dates are set
    var startDateStr = $('#activations-start-date').val();
    var endDateStr = $('#activations-end-date').val();

    if (startDateStr === '') {
      return;
    }

    if (endDateStr === '') {
      return;
    }

    // Check if start date is earlier as end date
    var startDate = new Date(startDateStr);
    var endDate = new Date(endDateStr);
    if (startDate >= endDate) {
      toastr.warning('End date has to be later than start date');
      return;
    }

    // TODO: Check user field

    // Assign start and
    $scope.activationsStartDate = startDateStr;
    $scope.activationsEndDate = endDateStr;

    $scope.activations = [];
    $scope.currentPage = 1;
    $scope.loadActivations();
  };

  var pickadateOptions = {
    format: 'yyyy-mm-dd',
    onClose: $scope.onFilterChange
  };
  $('.datepicker').pickadate(pickadateOptions);

  $scope.activations = [];
  $scope.currentPage = 1;
  $scope.itemsPerPage = 15;

  $scope.loadNextPage = function() {
    if ($scope.activations.length < $scope.itemsPerPage) {
      return;
    }

    $scope.currentPage++;
    $scope.loadActivations('', '', 0);
  };

  $scope.loadPrevPage = function() {
    if ($scope.currentPage < 1) {
      return;
    }

    $scope.currentPage--;
    $scope.loadActivations('', '', 0);
  };

  // Creates invoice on the server side and returns link
  $scope.exportSpreadsheet = function() {
    $http({
      method: 'POST',
      url: '/api/invoices',
      params: {
        startDate: $scope.activationsStartDate,
        endDate: $scope.activationsEndDate,
        //userId: 1,
        includeInvoiced: false,
        itemsPerPage: $scope.itemsPerPage,
        page: $scope.currentPage,
        ac: new Date().getTime()
      }
    })
    .success(function(invoiceData) {
      // invoiceData should contain link to the generated xls file
      toastr.success('Invoice created');
      console.log(invoiceData);

      var filePathParts = invoiceData.FilePath.split("/");
      var fileName = filePathParts[filePathParts.length-1];

      var alertContent = '<div class="row">'+
        '<div class="col-xs-6"><b>Invoice created!</b><br>'+
        '<b>File name:</b> ' + fileName + '</div>'+
        '<div class="col-xs-6"><a '+
        'href="/' + invoiceData.FilePath + '" '+ 
        'class="btn btn-primary btn-block">'+
        'Download</a></div>'+
        '</div>';
      vex.dialog.alert(alertContent);
    })
    .error(function() {
      toastr.error('Error creating invoice');
    });
  };

  $scope.deleteActivationPrompt = function(activationId) {
    var token = randomToken.generate();
    vex.dialog.prompt({
      message: 'Enter <span class="delete-prompt-token">' + 
       token + '</span> to delete',
      placeholder: 'Token',
      callback: $scope.deleteActivationPromptCallback.bind(this, token, activationId)
    });
  };

  $scope.deleteActivationPromptCallback = 
   function(expectedToken, activationId, value) {
    if (value) {    
      if (value === expectedToken) {
        $scope.deleteActivation(activationId);
      } else {
        toastr.error('Wrong token');
      }
    } else if (value !== false) {
      toastr.error('No token');
    }
  };

  $scope.deleteActivation = function(activationId) {
    $http({
      method: 'DELETE',
      url: '/api/activations/' + activationId,
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(data) {
      $scope.loadActivations();
    }) 
    .error(function() {
      toastr.error('Failed to delete activation');
    });
  };

  // Test Typeahead
  var substringMatcher = function(strs) {
    return function findMatches(q, cb) {
      var matches, substrRegex;
 
      // an array that will be populated with substring matches
      matches = [];
 
      // regex used to determine if a string contains the substring `q`
      substrRegex = new RegExp(q, 'i');
 
      // iterate through the pool of strings and for any string that
      // contains the substring `q`, add it to the `matches` array
      $.each(strs, function(i, str) {
        if (substrRegex.test(str)) {
          // the typeahead jQuery plugin expects suggestions to a
          // JavaScript object, refer to typeahead docs for more info
          matches.push({ value: str });
        }
      });
 
      cb(matches);
    };
  };
 
  var states = ['Alabama', 'Alaska', 'Arizona', 'Arkansas', 'California',
    'Colorado', 'Connecticut', 'Delaware', 'Florida', 'Georgia', 'Hawaii',
    'Idaho', 'Illinois', 'Indiana', 'Iowa', 'Kansas', 'Kentucky', 'Louisiana',
    'Maine', 'Maryland', 'Massachusetts', 'Michigan', 'Minnesota',
    'Mississippi', 'Missouri', 'Montana', 'Nebraska', 'Nevada', 'New Hampshire',
    'New Jersey', 'New Mexico', 'New York', 'North Carolina', 'North Dakota',
    'Ohio', 'Oklahoma', 'Oregon', 'Pennsylvania', 'Rhode Island',
    'South Carolina', 'South Dakota', 'Tennessee', 'Texas', 'Utah', 'Vermont',
    'Virginia', 'Washington', 'West Virginia', 'Wisconsin', 'Wyoming'
  ];
 
  $('.typeahead').typeahead({
    hint: true,
    highlight: true,
    minLength: 1
  },{
    name: 'states',
    displayKey: 'value',
    source: substringMatcher(states)
  });

}]); // app.controller

})(); // closure
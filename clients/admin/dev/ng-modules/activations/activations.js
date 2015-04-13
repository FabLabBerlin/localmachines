(function(){

'use strict';

var app = angular.module('fabsmith.admin.activations', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/activations', {
    templateUrl: 'ng-modules/activations/activations.html',
    controller: 'ActivationsCtrl'
  });
}]); // app.config

app.controller('ActivationsCtrl', ['$scope', '$http', '$location', 
 function($scope, $http, $location) {

  $('.datepicker').pickadate();
  $('.timepicker').pickatime();

  $scope.activations = [];
  $scope.currentPage = 1;
  $scope.itemsPerPage = 15;
  

  $scope.loadActivations = function(startDate, endDate, userId) {
    $http({
      method: 'GET',
      url: '/api/activations', 
      params: {
        startDate: '2015-01-01',
        endDate: '2015-02-01',
        userId: 1,
        includeInvoiced: false,
        itemsPerPage: $scope.itemsPerPage,
        page: $scope.currentPage
      }
    })
    .success(function(activations) {
      console.log(activations);
      $scope.activations = activations;
    })
    .error(function() {
      toastr.error('Failed to load activations');
    });
  };

  $scope.loadActivations();

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
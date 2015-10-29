(function(){

'use strict';

var app = angular.module('fabsmith.admin.dashboard', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
 $routeProvider.when('/dashboard', {
    templateUrl: 'ng-modules/dashboard/dashboard.html',
    controller: 'DashboardCtrl'
  });
}]); // app.config

app.controller('DashboardCtrl', ['$scope', '$http', '$location', 
 function($scope, $http, $location) {

  $scope.metrics = [];

  $scope.loadMetricsData = function() {
    $http({
      method: 'GET',
      url: '/api/metrics',
      params: {
        ac: new Date().getTime()
      }
    })
    .success(function(metrics) {
      $scope.metrics = metrics;
      $scope.renderChartsInit();
    })
    .error(function(data, status) {
      toastr.error('Failed to load metrics data');
    });
  };

  $scope.renderCharts = function(a, b) {
    var months = _.map($scope.metrics.ActivationsByMonth, function(sum, month) {
      return month;
    }).sort();
    var byMonth = months.map(function(month) {
      return [
        {
          v: month,
          f: month
        },
        Math.round($scope.metrics.ActivationsByMonth[month]),
        Math.round($scope.metrics.MembershipsByMonth[month])
      ];
    });


    var data = new google.visualization.DataTable();
    data.addColumn('string', 'Time of Day');
    data.addColumn('number', 'Activations (€)');
    data.addColumn('number', 'Memberships (€)');
    data.addRows(byMonth);

    var options = {
      title: 'Revenue through Activations and Memberships',
      hAxis: {
        title: 'Month',
      },
      vAxis: {
        title: 'Revenue / €'
      }
    };

    var chart = new google.visualization.ColumnChart(
      document.getElementById('chart_div'));

    chart.draw(data, options);
  };

  $scope.renderChartsInit = function() {
    $scope.renderCharts();
  };

  $scope.loadMetricsData();

}]); // app.controller

})(); // closure
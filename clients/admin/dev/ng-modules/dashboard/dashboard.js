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

  $scope.renderMonthlyCharts = function() {
    var months = _.map($scope.metrics.ActivationsByMonth, function(sum, month) {
      return month;
    }).sort();
    /*var byMonth = months.map(function(month) {
      var memberships = $scope.metrics.MembershipCountsByMonth[month];
      var minutes = Math.round($scope.metrics.MinutesByMonth[month]);
      return [
        {
          v: month,
          f: month + ' (' + memberships + ' non-free Memberships, ' + minutes + ' non-Admin minutes)'
        },
        Math.round($scope.metrics.ActivationsByMonth[month]),
        Math.round($scope.metrics.MembershipsByMonth[month])
      ];
    });*/


    /*var dataOld = new google.visualization.DataTable();
    dataOld.addColumn('string', 'Month');
    dataOld.addColumn('number', 'Activations (€)');
    dataOld.addColumn('number', 'Memberships (€)');
    dataOld.addRows(byMonth);*/

    var data = {
      labels: months,
      datasets: [
        {
          label: 'Activations (€)',
          fillColor: "rgba(220,220,220,0.5)",
          strokeColor: "rgba(220,220,220,0.8)",
          highlightFill: "rgba(220,220,220,0.75)",
          highlightStroke: "rgba(220,220,220,1)",
          data: _.map(months, function(month) {
            return Math.round($scope.metrics.ActivationsByMonth[month]);
          })
        },
        {
          label: 'Memberships (€)',
          fillColor: "rgba(151,187,205,0.5)",
          strokeColor: "rgba(151,187,205,0.8)",
          highlightFill: "rgba(151,187,205,0.75)",
          highlightStroke: "rgba(151,187,205,1)",
          data: _.map(months, function(month) {
            return Math.round($scope.metrics.MembershipsByMonth[month]);
          })
        }
      ]
    };

    var options = {
      title: 'Revenue through Activations and Memberships',
      hAxis: {
        title: 'Month',
      },
      vAxis: {
        title: 'Revenue / €'
      },
      legendTemplate : "<ul class=\"<%=name.toLowerCase()%>-legend\"><% for (var i=0; i<datasets.length; i++){%><li><span style=\"background-color:<%=datasets[i].strokeColor%>\"></span><%if(datasets[i].label){%><%=datasets[i].label%><%}%></li><%}%></ul>"
    };

    var ctx = document.getElementById('chart_monthly').getContext('2d');
    var chart = new Chart(ctx).Bar(data, options);
    /*var months = _.map($scope.metrics.ActivationsByMonth, function(sum, month) {
      return month;
    }).sort();
    var byMonth = months.map(function(month) {
      var memberships = $scope.metrics.MembershipCountsByMonth[month];
      var minutes = Math.round($scope.metrics.MinutesByMonth[month]);
      return [
        {
          v: month,
          f: month + ' (' + memberships + ' non-free Memberships, ' + minutes + ' non-Admin minutes)'
        },
        Math.round($scope.metrics.ActivationsByMonth[month]),
        Math.round($scope.metrics.MembershipsByMonth[month])
      ];
    });


    var data = new google.visualization.DataTable();
    data.addColumn('string', 'Month');
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
      document.getElementById('chart_monthly'));

    chart.draw(data, options);*/
  };

  $scope.renderDailyCharts = function() {
    var days = _.map($scope.metrics.ActivationsByDay, function(sum, day) {
      return day;
    }).sort();
    var byDay = days.map(function(day) {
      return [
        {
          v: moment(day).toDate(),
          f: day
        },
        Math.round($scope.metrics.ActivationsByDay[day])
      ];
    });


    var data = new google.visualization.DataTable();
    data.addColumn('date', 'Day');
    data.addColumn('number', 'Activations (€)');
    data.addRows(byDay);

    var options = {
      title: 'Daily Revenue through Activations',
      hAxis: {
        title: 'Day',
      },
      vAxis: {
        title: 'Revenue / €'
      }
    };

    var chart = new google.visualization.ColumnChart(
      document.getElementById('chart_daily'));

    chart.draw(data, options);
  };

  $scope.renderCharts = function() {
    $scope.renderMonthlyCharts();
    $scope.renderDailyCharts();
  };

  $scope.renderChartsInit = function() {
    $scope.renderCharts();
  };

  $scope.loadMetricsData();

}]); // app.controller

})(); // closure
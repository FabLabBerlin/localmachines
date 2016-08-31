(function(){

'use strict';

var app = angular.module('fabsmith.admin.dashboard', ['ngRoute', 'ngCookies']);

app.config(['$routeProvider', function($routeProvider) {
 $routeProvider.when('/dashboard', {
    templateUrl: 'ng-modules/dashboard/dashboard.html',
    controller: 'DashboardCtrl'
  });
}]); // app.config

app.controller('DashboardCtrl',
 ['$scope', '$http', '$location', '$cookies',
 function($scope, $http, $location, $cookies) {

  $scope.metrics = [];

  $scope.loadMetricsData = function() {
    $http({
      method: 'GET',
      url: '/api/metrics',
      params: {
        location: $cookies.get('location'),
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
    var byMonth = months.map(function(month) {
      var membershipsRnD = $scope.metrics.MembershipCountsByMonthRnD[month];
      var memberships = $scope.metrics.MembershipCountsByMonth[month] - membershipsRnD;
      var minutes = Math.round($scope.metrics.MinutesByMonth[month]);
      var activationsRevenue = Math.round($scope.metrics.ActivationsByMonth[month]);
      var membershipsRevenueRnd = Math.round($scope.metrics.MembershipsByMonthRnD[month] || 0);
      var membershipsRevenue = Math.round($scope.metrics.MembershipsByMonth[month]) - membershipsRevenueRnd;
      return [
        {
          v: month,
          f: month
        },
        membershipsRevenue,
        'Memberships (€): <b>' + membershipsRevenue + '</b><br>' + memberships + ' non-free Memberships',
        activationsRevenue,
        'Activations (€): <b>' + activationsRevenue + '</b><br>' + minutes + ' minutes for non-Admins',
        membershipsRevenueRnd,
        'R&D Center (€): <b>' + membershipsRevenueRnd + '</b><br>' + membershipsRnD + ' R&D Center Tables'
      ];
    });


    var data = new google.visualization.DataTable();
    data.addColumn('string', 'Month');
    data.addColumn('number', 'Memberships (€)');
    data.addColumn({'type': 'string', 'role': 'tooltip', 'p': {'html': true}});
    data.addColumn('number', 'Activations (€)');
    data.addColumn({'type': 'string', 'role': 'tooltip', 'p': {'html': true}});
    data.addColumn('number', 'R&D Center (€)');
    data.addColumn({'type': 'string', 'role': 'tooltip', 'p': {'html': true}});
    data.addRows(byMonth);

    var options = {
      title: 'Revenue through Activations and Memberships',
      hAxis: {
        title: 'Month',
      },
      vAxis: {
        title: 'Revenue / €'
      },
      tooltip: {isHtml: true},
      isStacked: true
    };

    var chart = new google.visualization.ColumnChart(
      document.getElementById('chart_monthly'));

    chart.draw(data, options);
  };

  $scope.renderDailyCharts = function() {
    var days = _.map($scope.metrics.ActivationsByDay, function(sum, day) {
      return day;
    }).sort();
    var byDay = days.map(function(day) {
      var minutes = Math.round($scope.metrics.MinutesByDay[day]);
      var activationsRevenue = Math.round($scope.metrics.ActivationsByDay[day]);
      return [
        {
          v: moment(day).toDate(),
          f: day
        },
        activationsRevenue,
        moment(day).format('D MMM YYYY') + '<br>Activations (€): <b>' + activationsRevenue + '</b><br>' + minutes + ' minutes for non-Admins'
      ];
    });


    var data = new google.visualization.DataTable();
    data.addColumn('date', 'Day');
    data.addColumn('number', 'Activations (€)');
    data.addColumn({'type': 'string', 'role': 'tooltip', 'p': {'html': true}});
    data.addRows(byDay);

    var options = {
      title: 'Daily Revenue through Activations',
      hAxis: {
        title: 'Day',
      },
      vAxis: {
        title: 'Revenue / €'
      },
      tooltip: {isHtml: true}
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

  var waitForFinalEvent = (function () {
    var timers = {};
    return function (callback, ms, uniqueId) {
      if (!uniqueId) {
        uniqueId = "Don't call this twice without a uniqueId";
      }
      if (timers[uniqueId]) {
        clearTimeout (timers[uniqueId]);
      }
      timers[uniqueId] = setTimeout(callback, ms);
    };
  })();

  $(window).resize(function(){
    waitForFinalEvent(function(){
      $scope.renderCharts();
    }, 500, "windowResize");
  });

  $scope.loadMetricsData();

}]); // app.controller

})(); // closure
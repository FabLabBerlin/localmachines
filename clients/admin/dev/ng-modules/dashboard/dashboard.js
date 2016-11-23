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
 ['$scope', '$http', '$location', '$cookies', 'api',
 function($scope, $http, $location, $cookies, api) {

  $scope.metrics = [];
  var currency = '';

  api.loadSettings(function(settings) {
    $scope.settings = settings;
    currency = $scope.settings.Currency.ValueString;
  });

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
      $http({
        method: 'GET',
        url: '/api/metrics/machine_earnings',
        params: {
          location: $cookies.get('location'),
        ac: new Date().getTime()
        }
      })
      .success(function(machineEarnings) {
        $scope.machineEarnings = machineEarnings;
        $scope.renderMachineEarnings();
      })
      .error(function(data, status) {
        toastr.error('Failed to load machine earnings data');
      });
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
        'Memberships (' + currency + '): <b>' + membershipsRevenue + '</b><br>' + memberships + ' non-free Memberships',
        activationsRevenue,
        'Activations (' + currency + '): <b>' + activationsRevenue + '</b><br>' + minutes + ' minutes for non-Admins',
        membershipsRevenueRnd,
        'R&D Center (' + currency + '): <b>' + membershipsRevenueRnd + '</b><br>' + membershipsRnD + ' R&D Center Tables'
      ];
    });


    var data = new google.visualization.DataTable();
    data.addColumn('string', 'Month');
    data.addColumn('number', 'Memberships (' + currency + ')');
    data.addColumn({'type': 'string', 'role': 'tooltip', 'p': {'html': true}});
    data.addColumn('number', 'Activations (' + currency + ')');
    data.addColumn({'type': 'string', 'role': 'tooltip', 'p': {'html': true}});
    data.addColumn('number', 'Co-Working (' + currency + ')');
    data.addColumn({'type': 'string', 'role': 'tooltip', 'p': {'html': true}});
    data.addRows(byMonth);

    var options = {
      title: 'Revenue through Activations and Memberships',
      hAxis: {
        title: 'Month',
      },
      vAxis: {
        title: 'Revenue / ' + currency
      },
      tooltip: {isHtml: true},
      isStacked: true,
      explorer: {
        axis: 'horizontal',
        keepInBounds: true
      }
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
        moment(day).format('D MMM YYYY') + '<br>Activations (' + currency + '): <b>' + activationsRevenue + '</b><br>' + minutes + ' minutes for non-Admins'
      ];
    });


    var data = new google.visualization.DataTable();
    data.addColumn('date', 'Day');
    data.addColumn('number', 'Activations (' + currency + ')');
    data.addColumn({'type': 'string', 'role': 'tooltip', 'p': {'html': true}});
    data.addRows(byDay);

    var options = {
      title: 'Daily Revenue through Activations',
      hAxis: {
        title: 'Day',
      },
      vAxis: {
        title: 'Revenue / ' + currency
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

  $scope.renderMachineEarnings = function() {
    console.log('$scope.machineEarnings=', $scope.machineEarnings);

    var ary = [
      ['Source', 'Memberships', 'Pay-As-You-Go', { role: 'annotation' } ]
    ];

    const typeNames = {
      0: 'z Other',
      1: '3D Printer',
      2: 'CNC mill',
      3: 'Heatpress',
      4: 'Knitting Machine',
      5: 'Lasercutters',
      6: 'Vinylcutter'
    };

    var sorted = _.sortBy($scope.machineEarnings, function(earning) {
      return [typeNames[earning.Machine.TypeId], earning.Machine.Name];
    });

    sorted = _.filter(sorted, function(earning) {
      return !earning.Machine.Archived;
    });

    _.each(sorted, function(earning) {
      ary.push([
        earning.Machine.Name, earning.Memberships, earning.PayAsYouGo, 0
      ]);
    });

    var data = google.visualization.arrayToDataTable(ary);

    var options = {
      title: 'Earnings by Machine',
      bars: 'horizontal',
      width: window.innerWidth,
      height: window.innerWidth,
      legend: { position: 'top', maxLines: 3 },
      //bar: { groupWidth: '75%' },
      isStacked: true,
      hAxis: {
        logscale: true
      }
    };

    var material = new google.charts.Bar(document.getElementById('chart_machine_earnings'));
    material.draw(data, options);
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
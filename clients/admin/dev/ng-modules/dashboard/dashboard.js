/*global metricsGcharts, metricsLoad */

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

  $scope.activeOnly = true;
  $scope.metrics = [];
  $scope.timeframe = {
    from: '2015-08-01',
    to: moment().format('YYYY-MM-DD')
  };
  $scope.binwidth = 'month';
  var currency = '';

  $(document).ready(function(){
    $('[data-toggle=tooltip]').hover(function(){
      // on mouseenter
      $(this).tooltip('show');
    }, function(){
      // on mouseleave
      $(this).tooltip('hide');
    });
  });

  api.loadSettings(function(settings) {
    $scope.settings = settings;
    currency = $scope.settings.Currency.ValueString;
  });

  $scope.loadMetricsData = function() {
    var locationId = $cookies.get('location');
    var options = {
      locationId: locationId,
      timeframe: $scope.timeframe,
      binwidth: $scope.binwidth
    };

    metricsLoad.main(options).then(function(metrics) {
      $scope.$apply(function() {
        $scope.metrics = metrics;
      });
      $scope.renderMainChart();
      metricsLoad.machineEarnings(options)
      .then(function(machineEarnings) {
        $scope.$apply(function() {
          $scope.machineEarnings = machineEarnings;
        });
        $scope.renderMachineEarnings();
        metricsLoad.machineCapacities(options)
        .then(function(machineCapacities) {
          $scope.$apply(function() {
            $scope.machineCapacities = machineCapacities;
          });
          $scope.renderMachineCapacities();
          metricsLoad.retention(options)
          .then(function(retention) {
            console.log('retention=', retention);
            $scope.$apply(function() {
              $scope.retention = retention;
              var pickadateOptions = {
                format: 'yyyy-mm-dd'
              };

              $('.datepicker').pickadate(pickadateOptions);
            });

            metricsLoad.memberships(options)
            .then(function(memberships) {
              console.log('memberships=', memberships);
              $scope.$apply(function() {
                $scope.memberships = memberships;
              });
              metricsGcharts.memberships(
                document.getElementById('chart_memberships'),
                memberships,
                currency
              );
              metricsLoad.heatmap(options)
              .then(function(coordinates) {
                $scope.$apply(function() {
                  $scope.coordinates = coordinates;
                });
                metricsGcharts.heatmap(
                  document.getElementById('heatmap_container'),
                  coordinates
                );
              });
            });
          });
        });
      });
    });
  };

  $scope.updateTimeframes = function() {
    $scope.metrics = undefined;
    $scope.machineCapacities = undefined;
    $scope.machineEarnings = undefined;
    $scope.retention = undefined;

    $scope.loadMetricsData();
  };

  $scope.renderMainChart = function() {
    console.log('$scope.metrics=', $scope.metrics);
    var hValues = _.map($scope.metrics.Activations, function(sum, hVal) {
      return hVal;
    }).sort();
    var zipped = hValues.map(function(hVal) {
      var membershipsRnD = $scope.metrics.MembershipCountsRnD[hVal];
      var memberships = $scope.metrics.MembershipCounts[hVal] - membershipsRnD;
      var minutes = Math.round($scope.metrics.Minutes[hVal]);
      var activationsRevenue = Math.round($scope.metrics.Activations[hVal]);
      var membershipsRevenueRnd = Math.round($scope.metrics.MembershipsRnD[hVal] || 0);
      var membershipsRevenue = Math.round($scope.metrics.Memberships[hVal]) - membershipsRevenueRnd;
      var title = '' + hVal + '<br><br>';
      return [
        {
          v: hVal,
          f: hVal
        },
        membershipsRevenue,
        title + 'Memberships (' + currency + '): <b>' + membershipsRevenue + '</b><br>' + memberships + ' non-free Memberships',
        activationsRevenue,
        title + 'Activations (' + currency + '): <b>' + activationsRevenue + '</b><br>' + minutes + ' minutes for non-Admins',
        membershipsRevenueRnd,
        title + 'R&D Center (' + currency + '): <b>' + membershipsRevenueRnd + '</b><br>' + membershipsRnD + ' R&D Center Tables'
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
    data.addRows(zipped);

    var options = {
      hAxis: {
        title: $scope.binwidth[0].toUpperCase() + $scope.binwidth.slice(1),
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
      document.getElementById('main_chart'));

    chart.draw(data, options);
  };

  $scope.renderChartsInit = function() {
    $scope.renderCharts();
  };

  $scope.renderMachineCapacities = function() {
    console.log('$scope.machineCapacities=', $scope.machineCapacities);

    var ary = [
      ['Source', 'Active (Days)', 'Lifetime (days)', { role: 'annotation' } ]
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

    var sorted = _.sortBy($scope.machineCapacities, function(c) {
      return [typeNames[c.Machine.TypeId], c.Machine.Name];
    });

    sorted = _.filter(sorted, function(c) {
      return !c.Machine.Archived;
    });

    _.each(sorted, function(c) {
      ary.push([
        c.Machine.Name, c.Hours, c.Capacity, 0
      ]);
    });

    var data = google.visualization.arrayToDataTable(ary);

    var options = {
      bars: 'horizontal',
      width: window.innerWidth * 0.9,
      height: window.innerWidth,
      legend: { position: 'top', maxLines: 3 },
      //bar: { groupWidth: '75%' },
      isStacked: true,
      hAxis: {
        logscale: true
      }
    };

    var material = new google.charts.Bar(document.getElementById('chart_machine_capacities'));
    material.draw(data, options);
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
      bars: 'horizontal',
      width: window.innerWidth * 0.9,
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

  $scope.formatPercentage = function(r) {
    return Math.round(100 * r);
  };

  $scope.retentionActiveClass = function(r) {
    return 'retention-' + Math.round(r / $scope.retention.activeMaxReturn * 4);
  };

  $scope.retentionClass = function(r) {
    return 'retention-' + Math.round(r / $scope.retention.maxReturn * 4);
  };

  $(window).resize(function(){
    waitForFinalEvent(function(){
      $scope.renderCharts();
    }, 500, "windowResize");
  });

  $scope.loadMetricsData();

}]); // app.controller

})(); // closure
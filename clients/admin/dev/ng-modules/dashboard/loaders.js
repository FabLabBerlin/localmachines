window.metricsLoad = {
  main: function(locId) {
    return $.ajax({
      method: 'GET',
      url: '/api/metrics',
      data: {
        location: locId
      }
    });
  },

  machineEarnings: function(locId) {
    return $.ajax({
      method: 'GET',
      url: '/api/metrics/machine_earnings',
      data: {
        location: locId
      }
    });
  },

  machineCapacities: function(locId) {
    return $.ajax({
      method: 'GET',
      url: '/api/metrics/machine_capacities',
      data: {
        location: locId
      }
    });
  },

  retention: function(locId) {
    var dfd = $.Deferred();

    $.when(
      $.ajax({
        method: 'GET',
        url: '/api/metrics/retention',
        data: {
          location: locId
        }
      }),
      $.ajax({
        method: 'GET',
        url: '/api/metrics/retention?excludeNeverActive=true',
        data: {
          location: locId
        }
      })
    ).done(function(r, rActive) {
      console.log('r=', r);
      console.log('rActive=', rActive);
      var retention = {
        all: r[0],
        active: rActive[0],
        maxReturn: undefined,
        activeMaxReturn: undefined
      };

      _.each(retention.all, function(row) {
        _.each(row.Returns, function(r) {
          if (_.isUndefined(retention.maxReturn) || retention.maxReturn < r) {
            retention.maxReturn = r;
          }
        });
      });
      console.log('retention.maxReturn=', retention.maxReturn);

      _.each(retention.active, function(row) {
        _.each(row.Returns, function(r) {
          if (_.isUndefined(retention.activeMaxReturn) || retention.activeMaxReturn < r) {
            retention.activeMaxReturn = r;
          }
        });
      });
      console.log('retention.maxReturn=', retention.activeMaxReturn);

      dfd.resolve(retention);
    }).fail(function() {
      dfd.reject('Error loading retention');
    });

    return dfd;
  }
};

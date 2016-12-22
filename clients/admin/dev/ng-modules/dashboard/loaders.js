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
    return $.ajax({
      method: 'GET',
      url: '/api/metrics/retention',
      data: {
        location: locId
      }
    });
  }
};

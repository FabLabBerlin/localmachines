// Material charts:
// https://github.com/google/google-visualization-issues/issues/2143
// var materialOptions = google.charts.Bar.convertOptions(classicOptions);


function membershipsCompactify(titles, titleByIdx, stats) {
  var total = 0;

  var revenueByTitle = {};
  _.each(stats, function(bin) {
    _.each(bin.UmsByName, function(ums) {
      _.each(ums, function(um) {
        if (!revenueByTitle[um.Membership.Title]) {
          revenueByTitle[um.Membership.Title] = 0;
        }

        revenueByTitle[um.Membership.Title] += um.Membership.MonthlyPrice;
        total += um.Membership.MonthlyPrice;
      });
    });
  });

  var filtered = [];
  _.each(revenueByTitle, function(revenue, title) {
    filtered.push({
      title: title,
      revenue: revenue
    });
  });

  filtered = _.sortBy(filtered, function(obj) {
    return -obj.revenue;
  });

  var splitAt = 10;

  var others = _.sortBy(_.pluck(filtered.slice(splitAt), 'title'), function(title) {
    return title;
  });
  filtered = _.pluck(filtered.slice(0, splitAt), 'title');
  filtered = _.sortBy(filtered, function(title) {
    return title;
  });

  titles = filtered;
  var filteredTitlesByIdx = {};
  _.each(titleByIdx, function(idx, title) {
    var i = _.indexOf(filtered, title);

    if (i >= 0) {
      filteredTitlesByIdx[title] = i;
    } else {
      filteredTitlesByIdx[title] = titles.length;
    }
  });

  titles.push('Other');

  $('#other-memberships').html('Other: ' +  others.join(', '));

  return {
    titles: titles,
    titleByIdx: filteredTitlesByIdx
  };
}


window.metricsGcharts = {
  memberships: function(domElement, stats, currency) {
    var titles = _.reduce(stats, function(res, bin) {
      _.each(bin.UmsByName, function(ums, title) {
        res.push(title);
      });
      return res;
    }, []);
    titles = _.uniq(titles);
    titles = _.sortBy(titles, function(title) {
      return title;
    });
    var titleByIdx = _.reduce(titles, function(res, title, i) {
      res[title] = i;
      return res;
    }, {});

    if (titles.length > 20) {
      var compact = membershipsCompactify(titles, titleByIdx, stats);
      titles = compact.titles;
      titleByIdx = compact.titleByIdx;
    }

    var data = [
      _.cloneDeep(titles)
    ];

    data[0].unshift('Membership');
    data[0].push({ role: 'annotation' });

    _.each(stats, function(bin) {
      var row = new Array(titles.length);
      row.unshift(bin.Month);

      _.each(bin.UmsByName, function(ums, title) {
        var i = titleByIdx[title];
        row[i + 1] = _.reduce(ums, function(sum, um) {
          return sum + um.Membership.MonthlyPrice;
        }, 0);
      });
      row.push('');

      data.push(row);
    });

    console.log('data=', data);

    data = google.visualization.arrayToDataTable(data);

    var options = {
      width: window.innerWidth,
      height: window.innerHeight,
      legend: { position: 'top', maxLines: 3 },
      bar: { groupWidth: '75%' },
      stacked: true,
      vAxis: {
        title: 'Revenue / ' + currency
      }
    };

    var material = new google.charts.Bar(domElement);
    material.draw(data, options);
  }
};

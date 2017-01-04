// Material charts:
// https://github.com/google/google-visualization-issues/issues/2143
// var materialOptions = google.charts.Bar.convertOptions(classicOptions);


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
    console.log('membershipTitles=', titles);
    console.log('titleByIdx=', titleByIdx);

    var data = [
      _.cloneDeep(titles)
    ];

    data[0].unshift('Membership');
    data[0].push({ role: 'annotation' });

    _.each(stats, function(bin) {
      var row = new Array(titles.length);
      console.log('bin.Month=', bin.Month);
      row.unshift(bin.Month);

      _.each(bin.UmsByName, function(ums, title) {
        var i = titleByIdx[title];
        console.log('ums=', ums);
        row[i + 1] = _.reduce(ums, function(sum, um) {
          return sum + um.Membership.MonthlyPrice;
        }, 0);
      });
      row.push('');

      data.push(row);
    });

    console.log('data=', data);

    /*var data = [
      ['Genre', 'Fantasy & Sci Fi', 'Romance', 'Mystery/Crime', 'General',
       'Western', 'Literature', { role: 'annotation' } ],
      ['2010', 10, 24, 20, 32, 18, 5, ''],
      ['2020', 16, 22, 23, 30, 16, 9, ''],
      ['2030', 28, 19, 29, 30, 12, 13, '']
    ];*/

    data = google.visualization.arrayToDataTable(data);

    var options = {
      width: 600,
      height: 900,
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

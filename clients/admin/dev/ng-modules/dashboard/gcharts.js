/*global OpenLayers, HeatCanvas */

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
  heatmap: function(domElement, coordinates) {
    var w = window.innerWidth * 0.9;
    var h = window.innerHeight;
    $(domElement).html('<div id="map" style="width:' + w + 'px; height:' + h + 'px;"></div>');
    var map, layer;
    map = new OpenLayers.Map ("map", {
      controls: [
          new OpenLayers.Control.Attribution(),
          new OpenLayers.Control.Navigation()
      ],
      maxExtent: new OpenLayers.Bounds(-20037508.34,-20037508.34,
                                       20037508.34,20037508.34),
      numZoomLevels: 12,
      maxResolution: 156543.0339,
      displayProjection: new OpenLayers.Projection("EPSG:4326"),
      units: 'm',
      projection: new OpenLayers.Projection("EPSG:4326")
    });
    
    var mapnik = new OpenLayers.Layer.OSM.Mapnik("Mapnik", {
       displayOutsideMaxExtent: true,
       wrapDateLine: true
    });
    map.addLayer(mapnik);
    map.setBaseLayer(mapnik);

    //var bounds = OpenLayers.Bounds.fromArray([70.4,15.2,136.2,53.7])
    //        .transform(map.displayProjection, map.getProjectionObject());
    //map.zoomToExtent(bounds);
    map.setCenter(new OpenLayers.LonLat(-20,0).transform(map.displayProjection, map.getProjectionObject()), 2);

    var size = map.getSize();
    if (size.h > 320) {
        map.addControl(new OpenLayers.Control.PanZoomBar());
    } else {
        map.addControl(new OpenLayers.Control.PanZoom());
    }
    
    var heatmap = new OpenLayers.Layer.HeatCanvas("Heat Canvas", map, {},
            {'step':0.5, 'degree':HeatCanvas.LINEAR, 'opacity':0.7});
    var data = _.map(coordinates, function(c) {
      heatmap.pushData(c.Coordinate.lat, c.Coordinate.lon, 10, c.UserId);
    });

    for(var i=0,l=data.length; i<l; i++) {
        
    }
    map.addLayer(heatmap);
    window.map = map;
    console.log('heatmap=', heatmap);
    //map.getView().on('change:resolution', function() {
    //  alert('1');
    //});
  },

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
      width: window.innerWidth * 0.9,
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

var _ = require('lodash');
var {formatDate, subtractVAT, toEuro, toCents} = require('./helpers');
var moment = require('moment');
var React = require('react');


function formatTime(t) {
  if (t) {
    var d = parseInt(t.toString(), 10);
    var h = Math.floor(d / 3600);
    var m = Math.floor(d % 3600 / 60);
    var s = Math.floor(d % 3600 % 60);
    var str = '';
    if (h) {
      str += String(h) + ' h ';
    }
    if (h || m) {
      str += String(m) + ' m ';
    }
    if (h || m || s) {
      str += String(s) + ' s ';
    }
    return str;
  }
}

var BillTables = React.createClass({
  render() {
    if (this.props.info.Activations && this.props.info.Activations.length !== 0) {
      var activationsByMonth = _.groupBy(this.props.info.Activations, function(info) {
        return moment(info.TimeStart).format('MMM YYYY');
      });
      var nodes = [];
      var lastMonth;
      var sumPriceInclVAT;
      var sumPriceExclVAT;
      var sumPriceVAT;
      var sumDurations;
      var i = 0;

      var getTotal = function() {
        return (
          <tr key={i++}>
            <td><label>Total</label></td>
            <td><label></label></td>
            <td><label>{formatTime(sumDurations)}</label></td>
            <td><label>{toEuro(sumPriceExclVAT)}</label> <i className="fa fa-eur"></i></td>
            <td><label>{toEuro(sumPriceVAT)}</label> <i className="fa fa-eur"></i></td>
            <td><label>{toEuro(sumPriceInclVAT)}</label> <i className="fa fa-eur"></i></td>
          </tr>
        );
      };

      this.props.info.Activations.reverse();

      _.each(this.props.info.Activations, function(info) {
        var month = moment(info.TimeStart).format('MMM YYYY');
        if (month !== lastMonth) {
          if (lastMonth) {
            nodes.push(getTotal());
          }
          nodes.push(<tr key={i++}><td colSpan={6}><h4>{month}</h4></td></tr>);
          nodes.push(
            <tr key={i++}>
              <th>Machine</th>
              <th>Date</th>
              <th>Time</th>
              <th>Price excl. VAT</th>
              <th>VAT (19%)</th>
              <th>Total</th>
            </tr>
          );
          sumPriceInclVAT = 0;
          sumPriceExclVAT = 0;
          sumPriceVAT = 0;
          sumDurations = 0;
        }
        var duration = moment.duration(moment(info.TimeEnd).diff(moment(info.TimeStart))).asSeconds();
        sumDurations += duration;
        var priceInclVAT = toCents(info.DiscountedTotal);
        var priceExclVAT = toCents(subtractVAT(info.DiscountedTotal));
        var priceVAT = priceInclVAT - priceExclVAT;
        sumPriceInclVAT += priceInclVAT;
        sumPriceExclVAT += priceExclVAT;
        sumPriceVAT += priceVAT;
        nodes.push(
          <tr key={i++}>
            <td>{info.MachineName}</td>
            <td>{formatDate(moment(info.TimeStart))}</td>
            <td>{formatTime(duration)}</td>
            <td>{toEuro(priceExclVAT)} <i className="fa fa-eur"></i></td>
            <td>{toEuro(priceVAT)} <i className="fa fa-eur"></i></td>
            <td>{toEuro(priceInclVAT)} <i className="fa fa-eur"></i></td>
          </tr>
        );
        lastMonth = month;
      });
      nodes.push(getTotal());
      return (
        <table className="table table-striped table-hover" >
          <thead></thead>
          <tbody>
            {nodes}
          </tbody>
        </table>
      );
    } else {
      return <p>You do not have any expenses.</p>;
    }
  }
});

export default BillTables;

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
    if (s) {
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
      return (
        <div>
          {_.map(activationsByMonth, function(activations, month) {
            return <BillTable key={month} activations={activations} month={month}/>;
          })}
        </div>
      );
    } else {
      return <p>You do not have any expenses.</p>;
    }
  }
});


var BillTable = React.createClass({
  render() {
    var BillNode;
    var sumPriceInclVAT = 0;
    var sumPriceExclVAT = 0;
    var sumPriceVAT = 0;
    var sumDurations = 0;

    if (this.props.activations && this.props.activations.length !== 0) {
      BillNode = this.props.activations.map(function(info, i) {
        var duration = moment.duration(moment(info.TimeEnd).diff(moment(info.TimeStart))).asSeconds();
        sumDurations += duration;
        var priceInclVAT = toCents(info.DiscountedTotal);
        var priceExclVAT = toCents(subtractVAT(info.DiscountedTotal));
        var priceVAT = priceInclVAT - priceExclVAT;
        sumPriceInclVAT += priceInclVAT;
        sumPriceExclVAT += priceExclVAT;
        sumPriceVAT += priceVAT;
        return (
          <tr key={i} >
            <td>{info.MachineName}</td>
            <td>{formatDate(moment(info.TimeStart))}</td>
            <td>{formatTime(duration)}</td>
            <td>{toEuro(priceExclVAT)} <i className="fa fa-eur"></i></td>
            <td>{toEuro(priceVAT)} <i className="fa fa-eur"></i></td>
            <td>{toEuro(priceInclVAT)} <i className="fa fa-eur"></i></td>
          </tr>
        );
      });
      return (
        <div>
          <h4>{this.props.month}</h4>
          <table className="table table-striped table-hover" >
            <thead>
              <tr>
                <th>Machine</th>
                <th>Date</th>
                <th>Time</th>
                <th>Price excl. VAT</th>
                <th>VAT (19%)</th>
                <th>Total</th>
              </tr>
            </thead>
            <tbody>
              {BillNode}
              <tr>
                <td><label>Total</label></td>
                <td><label></label></td>
                <td><label>{formatTime(sumDurations)}</label></td>
                <td><label>{toEuro(sumPriceExclVAT)}</label> <i className="fa fa-eur"></i></td>
                <td><label>{toEuro(sumPriceVAT)}</label> <i className="fa fa-eur"></i></td>
                <td><label>{toEuro(sumPriceInclVAT)}</label> <i className="fa fa-eur"></i></td>
              </tr>
            </tbody>
          </table>
        </div>
      );
    } else {
      return <p>You do not have any expenses for this month.</p>;
    }
  }
});

export default BillTables;

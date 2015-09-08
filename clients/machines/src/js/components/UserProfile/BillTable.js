var _ = require('lodash');
var getters = require('../../getters');
var moment = require('moment');
var React = require('react');
var reactor = require('../../reactor');
var {formatDate, subtractVAT, toEuro, toCents} = require('./helpers');


function formatDuration(t) {
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
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      monthlyBills: getters.getMonthlyBills
    };
  },

  render() {
    if (this.props.info && this.props.info.Activations && this.props.info.Activations.length !== 0) {

      var i = 0;
      var nodes = [];

      _.each(this.state.monthlyBills, function(bill) {
        if (i > 0) {
          nodes.push(<tr key={i++}><td colSpan={6}></td></tr>);
        }
        nodes.push(
          <tr key={i++}>
            <td colSpan={6}>
              <h4 className="text-left">{bill.month}</h4>
              <h5 className="text-left">({toEuro(bill.sums.total.priceInclVAT)} <i className="fa fa-eur"/> total incl. VAT)</h5>
            </td>
          </tr>
        );

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

        _.each(bill.activations, function(info) {
          nodes.push(
            <tr key={i++}>
              <td>{info.MachineName}</td>
              <td>{formatDate(info.TimeStart)}</td>
              <td>{formatDuration(info.duration)}</td>
              <td>{toEuro(info.priceExclVAT)} <i className="fa fa-eur"></i></td>
              <td>{toEuro(info.priceVAT)} <i className="fa fa-eur"></i></td>
              <td>{toEuro(info.priceInclVAT)} <i className="fa fa-eur"></i></td>
            </tr>
          );
        });

        nodes.push(
          <tr key={i++}>
            <td><label>Total Pay-As-You-Go</label></td>
            <td><label></label></td>
            <td><label>{formatDuration(bill.sums.activations.durations)}</label></td>
            <td><label>{toEuro(bill.sums.activations.priceExclVAT)}</label> <i className="fa fa-eur"></i></td>
            <td><label>{toEuro(bill.sums.activations.priceVAT)}</label> <i className="fa fa-eur"></i></td>
            <td><label>{toEuro(bill.sums.activations.priceInclVAT)}</label> <i className="fa fa-eur"></i></td>
          </tr>
        );
        nodes.push(
          <tr key={i++}>
            <td><label>Total Memberships</label></td>
            <td><label></label></td>
            <td><label></label></td>
            <td><label>{toEuro(bill.sums.memberships.priceExclVAT)}</label> <i className="fa fa-eur"></i></td>
            <td><label>{toEuro(bill.sums.memberships.priceVAT)}</label> <i className="fa fa-eur"></i></td>
            <td><label>{toEuro(bill.sums.memberships.priceInclVAT)}</label> <i className="fa fa-eur"></i></td>
          </tr>
        );
        nodes.push(
          <tr key={i++}>
            <td><label>Total</label></td>
            <td><label></label></td>
            <td><label></label></td>
            <td><label>{toEuro(bill.sums.total.priceExclVAT)}</label> <i className="fa fa-eur"></i></td>
            <td><label>{toEuro(bill.sums.total.priceVAT)}</label> <i className="fa fa-eur"></i></td>
            <td><label>{toEuro(bill.sums.total.priceInclVAT)}</label> <i className="fa fa-eur"></i></td>
          </tr>
        );
      });

      return (
        <table className="bill-table table table-striped table-hover" >
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

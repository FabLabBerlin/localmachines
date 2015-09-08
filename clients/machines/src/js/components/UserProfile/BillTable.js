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
  calculateMonthlyBills() {
    var monthlyBills = [];
    var activationsByMonth = _.groupBy(this.props.info.Activations, function(info) {
      return moment(info.TimeStart).format('MMM YYYY');
    });

    var membershipsByMonth = reactor.evaluateToJS(getters.getMembershipsByMonth);

    var months = this.months();
    _.each(months, function(month) {
      monthlyBills.push(this.calculateMonthlyBill(activationsByMonth, membershipsByMonth, month.format('MMM YYYY')));
    }.bind(this));
    monthlyBills.reverse();
    return monthlyBills;
  },

  calculateMonthlyBill(activationsByMonth, membershipsByMonth, month) {
    // All prices in Eurocent
    var monthlyBill = {
      month: month,
      activations: [],
      memberships: [],
      sums: {
        activations: {
          priceInclVAT: 0,
          priceExclVAT: 0,
          priceVAT: 0,
          durations: 0
        },
        memberships: {
          priceInclVAT: 0,
          priceExclVAT: 0,
          priceVAT: 0
        },
        total: {}
      }
    };

    _.each(activationsByMonth[month], function(info) {
      var duration = moment.duration(moment(info.TimeEnd).diff(moment(info.TimeStart))).asSeconds();
      monthlyBill.sums.durations += duration;
      var priceInclVAT = toCents(info.DiscountedTotal);
      var priceExclVAT = toCents(subtractVAT(info.DiscountedTotal));
      var priceVAT = priceInclVAT - priceExclVAT;
      monthlyBill.sums.activations.priceInclVAT += priceInclVAT;
      monthlyBill.sums.activations.priceExclVAT += priceExclVAT;
      monthlyBill.sums.activations.priceVAT += priceVAT;
      monthlyBill.activations.push({
        MachineName: info.MachineName,
        TimeStart: moment(info.TimeStart),
        duration: duration,
        priceExclVAT: priceExclVAT,
        priceVAT: priceVAT,
        priceInclVAT: priceInclVAT
      });
    });

    _.each(membershipsByMonth[month], function(membership) {
      var totalPrice = toCents(membership.MonthlyPrice);
      var priceExclVat = toCents(subtractVAT(membership.MonthlyPrice));
      var vat = totalPrice - priceExclVat;
      monthlyBill.memberships.push({
        startDate: moment(membership.StartDate),
        endDate: moment(membership.StartDate).add(membership.Duration, 'd'),
        priceExclVAT: priceExclVat,
        priceVAT: vat,
        priceInclVAT: totalPrice
      });
      monthlyBill.sums.memberships.priceInclVAT += totalPrice;
      monthlyBill.sums.memberships.priceExclVAT += priceExclVat;
      monthlyBill.sums.memberships.priceVAT += vat;
    });

    monthlyBill.sums.total = {
      priceInclVAT: monthlyBill.sums.activations.priceInclVAT + monthlyBill.sums.memberships.priceInclVAT,
      priceExclVAT: monthlyBill.sums.activations.priceExclVAT + monthlyBill.sums.memberships.priceExclVAT,
      priceVAT: monthlyBill.sums.activations.priceVAT + monthlyBill.sums.memberships.priceVAT
    };

    return monthlyBill;
  },

  months() {
    var userInfo = reactor.evaluateToJS(getters.getUserInfo);
    var months = [];
    var created = moment(userInfo.Created);
    if (created.unix() <= 0) {
      created = moment('2015-07-01');
    }
    for (var t = created; t.isBefore(moment()); t.add(1, 'd')) {
      months.push(t.clone());
    }
    months = _.uniq(months, function(month) {
      return month.format('MMM YYYY');
    });
    return months;
  },

  render() {
    if (this.props.info && this.props.info.Activations && this.props.info.Activations.length !== 0) {

      var i = 0;
      var nodes = [];

      this.props.info.Activations.reverse();

      _.each(this.calculateMonthlyBills(), function(bill) {
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

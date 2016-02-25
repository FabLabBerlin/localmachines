var _ = require('lodash');
var getters = require('../../getters');
var LoaderLocal = require('../LoaderLocal');
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
      str += String(h) + 'h ';
    }
    if (h || m) {
      str += String(m) + 'm ';
    }
    if (h || m || s) {
      str += String(s) + 's ';
    }
    return str;
  }
}

var BillTables = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      bill: getters.getBill,
      monthlyBills: getters.getMonthlyBills
    };
  },

  render() {
    if (this.state.bill) {

      var i = 0;
      var tables = [];

      _.each(this.state.monthlyBills, function(bill) {
        
        var caption = [];
        var thead = [];
        var tbody = [];
        var tfoot = [];

        caption.push( 
          <div key={i++}>
            <h4 className="text-left">{bill.month}</h4>
            <h5 className="text-left">
              ({toEuro(bill.sums.total.priceInclVAT)} 
              <i className="fa fa-eur"/> total incl. VAT)
            </h5>
          </div>
        );

        thead.push(
          <tr key={i++}>
            <th>Machine</th>
            <th>Date</th>
            <th>Time</th>
            <th>Price excl. VAT</th>
            <th>VAT (19%)</th>
            <th>Total</th>
          </tr>
        );

        _.each(bill.purchases, function(purchase) {
          var label = purchase.MachineName;
          switch (purchase.Type) {
          case 'activation':
            // already okay
            break;
          case 'co-working':
            label = 'Co-Working';
            break;
          case 'reservation':
            label += ' (Reservation)';
            break;
          case 'space':
            label = 'Space Booking';
            break;
          default:
            console.log('unhandled purchase type ', purchase.Type);
          }
          tbody.push(
            <tr key={i++}>
              <td>{label}</td>
              <td>{formatDate(purchase.TimeStart)}</td>
              <td>{formatDuration(purchase.duration)}</td>
              <td>{toEuro(purchase.priceExclVAT)}€</td>
              <td>{toEuro(purchase.priceVAT)}€</td>
              <td>{toEuro(purchase.priceInclVAT)}€</td>
            </tr>
          );
        });

        bill.purchases = _.sortBy(bill.purchases, (p) => {
          return -p.TimeStart.unix();
        });

        tfoot.push(
          <tr key={i++}>
            <td><b>Total Pay-As-You-Go</b></td>
            <td>&nbsp;</td>
            <td><b>{formatDuration(bill.sums.purchases.durations)}</b></td>
            <td><b>{toEuro(bill.sums.purchases.priceExclVAT)}€</b></td>
            <td><b>{toEuro(bill.sums.purchases.priceVAT)}€</b></td>
            <td><b>{toEuro(bill.sums.purchases.priceInclVAT)}€</b></td>
          </tr>
        );

        tfoot.push(
          <tr key={i++}>
            <td><b>Total Memberships</b></td>
            <td>&nbsp;</td>
            <td>&nbsp;</td>
            <td><b>{toEuro(bill.sums.memberships.priceExclVAT)}€</b></td>
            <td><b>{toEuro(bill.sums.memberships.priceVAT)}€</b></td>
            <td><b>{toEuro(bill.sums.memberships.priceInclVAT)}€</b></td>
          </tr>
        );

        tfoot.push(
          <tr key={i++}>
            <td><b>Total</b></td>
            <td>&nbsp;</td>
            <td>&nbsp;</td>
            <td><b>{toEuro(bill.sums.total.priceExclVAT)}€</b></td>
            <td><b>{toEuro(bill.sums.total.priceVAT)}€</b></td>
            <td><b>{toEuro(bill.sums.total.priceInclVAT)}€</b></td>
          </tr>
        );

        tables.push(
          <div key={i++}>
            {caption}
            <div className="table-responsive">
              <table className="table table-stripped table-hover">
                <thead>{thead}</thead>
                <tbody>{tbody}</tbody>
                <tfoot>{tfoot}</tfoot>
              </table>
            </div>
          </div>
        );
      });

      return (
        <div>
          {tables}
        </div>
      );
    } else {
      return <LoaderLocal/>;
    }
  }
});

export default BillTables;

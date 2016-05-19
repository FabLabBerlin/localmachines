var _ = require('lodash');
var moment = require('moment');
var React = require('react');
var reactor = require('../../reactor');
var SettingsGetters = require('../../modules/Settings/getters');
var {formatDate, subtractVAT, toEuro, toCents} = require('./helpers');

function formatDuration(purchase) {
  if (purchase.Quantity) {
    var duration = purchase.Quantity;
    switch (purchase.PriceUnit) {
    case 'month':
      duration *= 60 * 60 * 24 * 30;
      break;
    case 'day':
      duration *= 60 * 60 * 24;
      break;
    case 'hour':
      duration *= 60 * 60;
      break;
    case '30 minutes':
      duration *= 60 * 30;
      break;
    case 'minute':
      duration *= 60;
      break;
    case 'second':
      break;
    default:
      console.log('unknown price unit', purchase.PriceUnit);
      return undefined;
    }

    var d = parseInt(duration.toString(), 10);
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

function formatPrice(price) {
  return (Math.round(price * 100) / 100).toFixed(2);
}

var BillTable = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      vatPercent: SettingsGetters.getVatPercent
    };
  },

  render() {
    const bill = this.props.bill;

    var i = 0;
    var caption = [];
    var thead = [];
    var tbody = [];
    var tfoot = [];

    caption.push( 
      <div key={i++}>
        <h4 className="text-left">{bill.month}</h4>
        <h5 className="text-left">
          ({formatPrice(bill.Sums.All.PriceInclVAT)} 
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
        <th>VAT ({this.state.vatPercent}%)</th>
        <th>Total</th>
      </tr>
    );

    _.each(bill.Purchases.Data, function(purchase) {
      var label = purchase.Machine ? purchase.Machine.Name : '';
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
          <td>{formatDate(moment(purchase.TimeStart))}</td>
          <td>{formatDuration(purchase)}</td>
          <td>{formatPrice(purchase.PriceExclVAT)}€</td>
          <td>{formatPrice(purchase.PriceVAT)}€</td>
          <td>{formatPrice(purchase.DiscountedTotal)}€</td>
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
        <td></td>
        <td><b>{formatPrice(bill.Sums.Purchases.PriceExclVAT)}€</b></td>
        <td><b>{formatPrice(bill.Sums.Purchases.PriceVAT)}€</b></td>
        <td><b>{formatPrice(bill.Sums.Purchases.PriceInclVAT)}€</b></td>
      </tr>
    );

    tfoot.push(
      <tr key={i++}>
        <td><b>Total Memberships</b></td>
        <td>&nbsp;</td>
        <td>&nbsp;</td>
        <td><b>{formatPrice(bill.Sums.Memberships.PriceExclVAT)}€</b></td>
        <td><b>{formatPrice(bill.Sums.Memberships.PriceVAT)}€</b></td>
        <td><b>{formatPrice(bill.Sums.Memberships.PriceInclVAT)}€</b></td>
      </tr>
    );

    tfoot.push(
      <tr key={i++}>
        <td><b>Total</b></td>
        <td>&nbsp;</td>
        <td>&nbsp;</td>
        <td><b>{formatPrice(bill.Sums.All.PriceExclVAT)}€</b></td>
        <td><b>{formatPrice(bill.Sums.All.PriceVAT)}€</b></td>
        <td><b>{formatPrice(bill.Sums.All.PriceInclVAT)}€</b></td>
      </tr>
    );

    return (
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
  }
});

export default BillTable;

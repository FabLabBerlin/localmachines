var _ = require('lodash');
var $ = require('jquery');
var Invoices = require('../../modules/Invoices');
var LocationGetters = require('../../modules/Location/getters');
var moment = require('moment');
var React = require('react');
var reactor = require('../../reactor');
var SettingsGetters = require('../../modules/Settings/getters');
var toastr = require('../../toastr');
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
    var h = String(Math.floor(d / 3600));
    var m = String(Math.floor(d % 3600 / 60));
    var s = String(Math.floor(d % 3600 % 60));
    if (h.length === 1) {
      h = '0' + h;
    }
    if (m.length === 1) {
      m = '0' + m;
    }
    if (s.length === 1) {
      s = '0' + s;
    }
    var str = h + ':' + m + ':' + s + ' h';

    return str;
  }
}

function formatPrice(price) {
  return (Math.round(price * 100) / 100).toFixed(2);
}


var DurationEdit = React.createClass({
  componentDidMount() {
    $(this.refs.duration.getDOMNode()).focus()
                                      .select();
  },

  render() {
    return (
      <input type="text"
             ref="duration"
             value={formatDuration(this.props.purchase)}/>
    );
  }
});



var BillTable = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      editPurchaseId: Invoices.getters.getEditPurchaseId,
      isAdmin: LocationGetters.getIsAdmin,
      vatPercent: SettingsGetters.getVatPercent
    };
  },

  edit(purchase, e) {
    if (this.state.isAdmin) {
      Invoices.actions.editPurchase(purchase.Id);
    }
    e.stopPropagation();
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
        <h4 className="text-left">{bill.Month} / {bill.Year}</h4>
        <h5 className="text-left">
          ({formatPrice(bill.Sums.All.PriceInclVAT)} 
          <i className="fa fa-eur"/> total incl. VAT)
        </h5>
      </div>
    );

    tbody.push(
      <tr key={i++}>
        <td><b>Purchases</b></td>
        <td>&nbsp;</td>
        <td></td>
        <td></td>
        <td></td>
        <td></td>
      </tr>
    );

    tbody.push(
      <tr key={i++}>
        <th>Machine</th>
        <th>Date</th>
        <th>Time</th>
        <th>Price excl. VAT</th>
        <th>VAT ({this.state.vatPercent}%)</th>
        <th>Total</th>
      </tr>
    );

    var sorted = _.sortBy(bill.Purchases, (p) => {
      return -moment(p.TimeStart).unix();
    });

    _.each(sorted, (purchase) => {
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
      const selected = this.state.editPurchaseId === purchase.Id;
      tbody.push(
        <tr key={i++}
            onClick={this.edit.bind(this, purchase)}
            className={!selected ? 'unselected' : undefined}>
          <td>{label}</td>
          <td>{formatDate(moment(purchase.TimeStart))}</td>
          <td>
            {selected ?
              <DurationEdit purchase={purchase}/> :
              formatDuration(purchase)
            }
          </td>
          <td>{formatPrice(purchase.PriceExclVAT)}€</td>
          <td>{formatPrice(purchase.PriceVAT)}€</td>
          <td>{formatPrice(purchase.DiscountedTotal)}€</td>
        </tr>
      );
    });

    tbody.push(
      <tr key={i++}>
        <td><b>Memberships</b></td>
        <td>&nbsp;</td>
        <td></td>
        <td></td>
        <td></td>
        <td></td>
      </tr>
    );

    tbody.push(
      <tr key={i++}>
        <th>Type</th>
        <th>Start Date</th>
        <th>End Date</th>
        <th></th>
        <th></th>
        <th>Total</th>
      </tr>
    );

    _.each(bill.UserMemberships.Data, (um) => {
      tbody.push(
        <tr key={i++}>
          <td>{um.Title}</td>
          <td>{formatDate(moment(um.StartDate))}</td>
          <td>{formatDate(moment(um.EndDate))}</td>
          <td></td>
          <td></td>
          <td>{um.Bill ? (formatPrice(um.MonthlyPrice) + '€') : '-'}</td>
        </tr>
      );
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

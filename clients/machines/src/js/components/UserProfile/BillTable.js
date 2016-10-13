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
  render() {
    return (
      <input type="text"
             autoFocus="on"
             ref="duration"
             onChange={this.update}
             value={this.props.purchase.editedDuration ||
                    formatDuration(this.props.purchase)}/>
    );
  },

  update(e) {
    Invoices.actions.editPurchaseDuration(e.target.value);
  }
});


var EmptyRow = React.createClass({
  render() {
    const style = {
      border: 'none'
    };

    return (
      <tr>
        <td style={style}>
          {this.props.children}
        </td>
      </tr>
    );
  }
});


var AddPurchase = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      isAdmin: LocationGetters.getIsAdmin
    };
  },

  render() {
    if (this.state.isAdmin) {
      return (
        <button id="inv-add-purchase">
          <div id="inv-add-purchase-icon"/>
          <div>Add Purchase</div>
        </button>
      );
    } else {
      return <div/>;
    }
  }
});


var BillTable = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      currency: SettingsGetters.getCurrency,
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
    if (!this.props.bill) {
      return <div/>;
    }

    const bill = this.props.bill.toJS();

    var i = 0;
    var caption = [];
    var thead = [];
    var tbody = [];
    var tfoot = [];

    caption.push( 
      <div key={i++}>
        <h4 className="text-left">{moment().month(bill.Month - 1).format('MMMM')} / {bill.Year}</h4>
        <h5 className="text-left">
          ({formatPrice(bill.Sums.All.PriceInclVAT)} {this.state.currency} total incl. VAT)
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
        <th>Duration</th>
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
      case 'reservation':
        label += ' (Reservation)';
        break;
      case 'tutor':
        label = 'Tutoring';
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
          <td>{formatDate(moment(purchase.TimeStart))} {moment(purchase.TimeStart).format('HH:mm')}</td>
          <td>
            {selected && moment(purchase.TimeEnd).unix() > 0 ?
              <DurationEdit purchase={purchase}/> :
              (purchase.editedDuration ? purchase.editedDuration :
              formatDuration(purchase))
            }
          </td>
          <td>{formatPrice(purchase.PriceExclVAT)} {this.state.currency}</td>
          <td>{formatPrice(purchase.PriceVAT)} {this.state.currency}</td>
          <td>{formatPrice(purchase.DiscountedTotal)} {this.state.currency}</td>
        </tr>
      );
    });

    tbody.push(
      <EmptyRow key={i++}>
        <AddPurchase/>
      </EmptyRow>
    );
    tbody.push(<EmptyRow key={i++}/>);
    tbody.push(<EmptyRow key={i++}/>);
    tbody.push(<EmptyRow key={i++}/>);
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
          <td>{um.MonthlyPrice ? (formatPrice(um.MonthlyPrice) + this.state.currency) : '-'}</td>
        </tr>
      );
    });

    tfoot.push(
      <tr key={i++}>
        <td><b>Total Pay-As-You-Go</b></td>
        <td>&nbsp;</td>
        <td></td>
        <td><b>{formatPrice(bill.Sums.Purchases.PriceExclVAT)} {this.state.currency}</b></td>
        <td><b>{formatPrice(bill.Sums.Purchases.PriceVAT)} {this.state.currency}</b></td>
        <td><b>{formatPrice(bill.Sums.Purchases.PriceInclVAT)} {this.state.currency}</b></td>
      </tr>
    );

    tfoot.push(
      <tr key={i++}>
        <td><b>Total Memberships</b></td>
        <td>&nbsp;</td>
        <td>&nbsp;</td>
        <td><b>{formatPrice(bill.Sums.Memberships.PriceExclVAT)} {this.state.currency}</b></td>
        <td><b>{formatPrice(bill.Sums.Memberships.PriceVAT)} {this.state.currency}</b></td>
        <td><b>{formatPrice(bill.Sums.Memberships.PriceInclVAT)} {this.state.currency}</b></td>
      </tr>
    );

    tfoot.push(
      <tr key={i++}>
        <td><b>Total</b></td>
        <td>&nbsp;</td>
        <td>&nbsp;</td>
        <td><b>{formatPrice(bill.Sums.All.PriceExclVAT)} {this.state.currency}</b></td>
        <td><b>{formatPrice(bill.Sums.All.PriceVAT)} {this.state.currency}</b></td>
        <td><b>{formatPrice(bill.Sums.All.PriceInclVAT)} {this.state.currency}</b></td>
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

var _ = require('lodash');
var $ = require('jquery');
var Button = require('../Button');
var Edit = require('./PurchaseEditing');
var Invoices = require('../../modules/Invoices');
var LocationGetters = require('../../modules/Location/getters');
var moment = require('moment');
var React = require('react');
var reactor = require('../../reactor');
var SettingsGetters = require('../../modules/Settings/getters');
var toastr = require('../../toastr');
var {formatDate, formatDuration, formatPrice, subtractVAT, toEuro, toCents} = require('./helpers');

// https://github.com/HubSpot/vex/issues/72
var vex = require('vex-js'),
VexDialog = require('vex-js/js/vex.dialog.js');
vex.defaultOptions.className = 'vex-theme-custom';


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

  add() {
    Invoices.actions.createPurchase(this.props.invoice);
  },

  render() {
    if (this.state.isAdmin && this.props.visible) {
      return <Button.Annotated id="inv-add-purchase"
                               icon="/machines/assets/img/invoicing/add_purchase.svg"
                               label="Add Purchase"
                               onClick={this.add}/>;
    } else {
      return <div/>;
    }
  }
});


var RemovePurchase = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      isAdmin: LocationGetters.getIsAdmin
    };
  },

  archive() {
    VexDialog.buttons.YES.text = 'Yes';
    VexDialog.buttons.NO.text = 'No';

    VexDialog.confirm({
      message: 'Do you really want to archive this purchase?',
      callback: confirmed => {
        if (confirmed) {
          Invoices.actions.archivePurchase({
            invoice: this.props.invoice,
            purchaseId: this.props.purchase.Id
          });
        }
        $('.vex').remove();
        $('body').removeClass('vex-open');
      }
    });
  },

  render() {
    if (this.state.isAdmin && this.props.visible) {
      return (
        <button id="inv-remove-purchase"
                onClick={this.archive}>
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
    if (!this.props.invoice) {
      return <div/>;
    }

    const bill = this.props.invoice.toJS();

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
        <th>Category</th>
        <th>Name</th>
        <th>Date</th>
        <th>Amount</th>
        <th>Unit</th>
        <th>Price / Unit</th>
        <th>Total (incl. VAT)</th>
      </tr>
    );

    var sorted = _.sortBy(bill.Purchases, (p) => {
      return -moment(p.TimeStart).unix();
    });

    _.each(sorted, (purchase) => {
      var category;
      switch (purchase.Type) {
      case 'activation':
        category = 'Activation';
        break;
      case 'other':
        category = 'Other';
        break;
      case 'reservation':
        category = 'Reservation';
        break;
      case 'tutor':
        category = 'Tutoring';
        break;
      default:
        console.log('unhandled purchase type ', purchase.Type);
      }
      const selected = this.state.editPurchaseId === purchase.Id && this.props.addPurchaseVisible;
      const editable = selected && moment(purchase.TimeEnd).unix() > 0;

      tbody.push(
        <tr key={i++}
            onClick={this.edit.bind(this, purchase)}
            className={'inv-purchase ' + (selected ? 'selected' : 'unselected')}>
          <td>
            {editable ?
              <Edit.Category invoice={this.props.invoice} purchase={purchase}/> :
              category
            }
          </td>
          <td>
            {editable ? (
              <Edit.Name invoice={this.props.invoice} purchase={purchase}/>
            ) : (
              purchase.Type !== 'other' ?
              (purchase.Machine && purchase.Machine.Name) :
              purchase.CustomName
            )}
          </td>
          <td>
            {editable ? (
              <Edit.StartTime invoice={this.props.invoice} purchase={purchase}/>
            ) : (
              formatDate(moment(purchase.TimeStart)) + ' ' + moment(purchase.TimeStart).format('HH:mm')
            )}
          </td>
          <td>
            {editable ?
              <Edit.Amount invoice={this.props.invoice} purchase={purchase}/> :
              (purchase.editedDuration ? purchase.editedDuration : (
              purchase.PriceUnit !== 'gram' ? formatDuration(purchase) :
              purchase.Quantity))
            }
          </td>
          <td>
            {editable ?
              <Edit.Unit invoice={this.props.invoice} purchase={purchase}/> :
              purchase.PriceUnit
            }
          </td>
          <td>
            {editable ?
              <Edit.PricePerUnit invoice={this.props.invoice} purchase={purchase}/> :
              purchase.PricePerUnit
            }
          </td>
          <td>{formatPrice(purchase.DiscountedTotal)} {this.state.currency}</td>
          <td>
            {editable ?
              <RemovePurchase invoice={this.props.invoice}
                              purchase={purchase}
                              visible={this.props.addPurchaseVisible}/> :
              null
            }
          </td>
        </tr>
      );
    });

    tbody.push(
      <EmptyRow key={i++}>
        <AddPurchase invoice={this.props.invoice}
                     visible={this.props.addPurchaseVisible}/>
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
          <table className="table table-stripped">
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

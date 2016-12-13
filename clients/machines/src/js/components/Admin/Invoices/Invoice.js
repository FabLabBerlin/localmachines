var BillTable = require('../../UserProfile/BillTable');
var getters = require('../../../getters');
var Invoices = require('../../../modules/Invoices');
var LoaderLocal = require('../../LoaderLocal');
var Location = require('../../../modules/Location');
var moment = require('moment');
var React = require('react');
var reactor = require('../../../reactor');
var SettingsActions = require('../../../modules/Settings/actions');
var SettingsGetters = require('../../../modules/Settings/getters');
var toastr = require('../../../toastr');
var UserActions = require('../../../actions/UserActions');
var util = require('./util');

import {hashHistory} from 'react-router';


var Header = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      currency: SettingsGetters.getCurrency,
      editPurchaseId: Invoices.getters.getEditPurchaseId,
      invoicesActions: Invoices.getters.getInvoicesActions,
      invoiceStatuses: Invoices.getters.getInvoiceStatuses,
      location: Location.getters.getLocation,
      locationId: Location.getters.getLocationId,
      MonthlySums: Invoices.getters.getMonthlySums,
      uid: getters.getUid,
      userMemberships: Invoices.getters.getUserMemberships,
      vatPercent: SettingsGetters.getVatPercent
    };
  },

  componentWillMount() {
    const locationId = reactor.evaluateToJS(Location.getters.getLocationId);
    SettingsActions.loadSettings({locationId});
  },

  cancel(e) {
    e.stopPropagation();

    Invoices.actions.cancel(this.props.invoice);
  },

  complete(e) {
    e.stopPropagation();
    const invoice = this.props.invoice;

    Invoices.actions.complete(invoice);
  },

  hide() {
    if (this.state.editPurchaseId) {
      /*eslint-disable no-alert */
      if (window.confirm('Discard unsaved changes?')) {
        Invoices.actions.editPurchase(null);
      } else {
        return;
      }
      /*eslint-enable no-alert */
    }
    hashHistory.push('/admin/invoices');
  },

  makeDraft(e) {
    e.stopPropagation();
    const locId = this.state.locationId;
    const invoice = this.props.invoice;

    Invoices.actions.makeDraft(locId, invoice);
  },

  render() {
    const invoice = this.props.invoice;
    const invoiceStatuses = this.state.invoiceStatuses;

    if (!invoice) {
      return <LoaderLocal/>;
    }

    const amount = (Math.round(invoice.getIn(['Sums', 'All', 'PriceInclVAT']) * 100)
                    / 100).toFixed(2);
    const name = invoice.getIn(['User', 'FirstName']) +
                 ' ' + invoice.getIn(['User', 'LastName']);
    const timeFrame = '' + moment().month(invoice.get('Month') - 1).format('MMM') + '/' +
                           invoice.get('Year');
    const invoiceActions = this.state.invoicesActions.get(invoice.get('Id'));

    return (
      <div id="inv-header">
        <div className="row">
          <div className="col-sm-3 col-sm-push-9" style={{overflow: 'hidden'}}>
            <button type="button"
                    title="Close"
                    onClick={this.hide}>
              <i className="fa fa-close"/>
            </button>
            {invoiceActions.get('PushDraft') ?
              <button type="button"
                      onClick={this.makeDraft}
                      title="Upload Invoice Draft to Fastbill">
                <i className="fa fa-cloud-upload"/>
              </button> : null
            }
            {invoiceActions.get('Save') ?
              <button type="button"
                      onClick={this.save}
                      title="Save">
                <i className="fa fa-save"/>
              </button> : null
            }
            {invoiceActions.get('Freeze') ?
              <button type="button"
                      onClick={this.complete}
                      title="Create Invoice">
                <i className="fa fa-money"/>
              </button> : null
            }
            {invoiceActions.get('Send') ?
              <button type="button"
                      onClick={this.send}
                      title="Send Invoice">
                <img src="/machines/assets/img/invoicing/send_invoice_white.svg"/>
              </button> : null
            }
            {invoiceActions.get('SendCanceled') ?
              <button type="button"
                      onClick={this.sendCanceled}
                      title="Send Cancelation">
                <img src="/machines/assets/img/invoicing/send_invoice_white.svg"/>
              </button> : null
            }
            {invoiceActions.get('Cancel') ?
              <button type="button"
                      onClick={this.cancel}
                      title="Cancel Invoice">
                <i className="fa fa-ban"/>
              </button> : null
            }
          </div>
          <div className="col-sm-3 col-sm-pull-3">
            <h3>{name}</h3>
          </div>
          <div className="col-sm-3 col-sm-pull-3">
            <h3>
              {invoice.getIn(['User', 'NoAutoInvoicing']) ?
                'Manual Invoicing' : 'Automatic Invoicing'}
            </h3>
          </div>
          <div className="col-sm-3 col-sm-pull-3 inv-amount">
            <h3>Sum: {amount} {this.state.currency}</h3>
          </div>
        </div>
        <div className="row">
          <div className="col-sm-3">
            Fastbill Customer No: {invoice.getIn(['User', 'ClientId'])}
          </div>
          <div className="col-sm-3">
            {this.statusInfo()}
          </div>
          <div className="col-sm-3">
            Incl. {this.state.vatPercent}% VAT
          </div>
          <div className="col-sm-3">
          </div>
        </div>
      </div>
    );
  },

  save(e) {
    e.stopPropagation();
    const locId = this.state.locationId;
    const invoice = this.props.invoice;

    Invoices.actions.save(locId, {invoice});
  },

  send(e) {
    e.stopPropagation();

    Invoices.actions.send(false, this.props.invoice);
  },

  sendCanceled(e) {
    e.stopPropagation();

    Invoices.actions.sendCanceled(this.props.invoice);
  },

  statusInfo() {
    const invoice = this.props.invoice;

    return util.statusInfo(invoice);
  }
});


var Invoice = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      editPurchaseId: Invoices.getters.getEditPurchaseId,
      invoices: Invoices.getters.getInvoices,
      invoicesActions: Invoices.getters.getInvoicesActions,
      invoiceStatuses: Invoices.getters.getInvoiceStatuses,
      location: Location.getters.getLocation,
      locationId: Location.getters.getLocationId,
      MonthlySums: Invoices.getters.getMonthlySums,
      uid: getters.getUid,
      userMemberships: Invoices.getters.getUserMemberships
    };
  },

  componentWillMount() {
    const uid = reactor.evaluateToJS(getters.getUid);
    const invoiceId = parseInt(this.props.params.invoiceId);

    UserActions.fetchUser(uid);
    Location.actions.loadUserLocations(uid);
    Invoices.actions.fetchInvoice(this.state.locationId, {
      invoiceId: invoiceId
    });
  },

  render() {
    const invoiceId = parseInt(this.props.params.invoiceId);
    const invoice = this.state.invoices.get(invoiceId);
    const invoiceStatuses = this.state.invoiceStatuses;

    if (invoice) {
      return (
        <div className="inv-monthly-sums">
          <div className="inv-invoice"
               onClick={this.stopPropagation}>
            <Header invoice={invoice}/>
            <div id="inv-body">
              <BillTable invoice={invoice}
                         addPurchaseVisible={true}/>
            </div>
          </div>
        </div>
      );
    } else {
      return <LoaderLocal/>;
    }
  },

  stopPropagation(e) {
    e.stopPropagation();
  }

});

export default Invoice;

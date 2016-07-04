var BillTable = require('../../UserProfile/BillTable');
var getters = require('../../../getters');
var Invoices = require('../../../modules/Invoices');
var LoaderLocal = require('../../LoaderLocal');
var LocationGetters = require('../../../modules/Location/getters');
var moment = require('moment');
var React = require('react');
var reactor = require('../../../reactor');
var SettingsActions = require('../../../modules/Settings/actions');
var SettingsGetters = require('../../../modules/Settings/getters');
var toastr = require('../../../toastr');


var Header = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      editPurchaseId: Invoices.getters.getEditPurchaseId,
      invoice: Invoices.getters.getInvoice,
      invoiceActions: Invoices.getters.getInvoiceActions,
      invoiceStatuses: Invoices.getters.getInvoiceStatuses,
      location: LocationGetters.getLocation,
      locationId: LocationGetters.getLocationId,
      MonthlySums: Invoices.getters.getMonthlySums,
      uid: getters.getUid,
      userMemberships: Invoices.getters.getUserMemberships,
      vatPercent: SettingsGetters.getVatPercent
    };
  },

  componentWillMount() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    SettingsActions.loadSettings({locationId});
  },

  cancel(e) {
    e.stopPropagation();

    console.log('Invoice#cancel');

    Invoices.actions.cancel();
  },

  complete(e) {
    e.stopPropagation();
    const locId = this.state.locationId;
    const invoice = this.state.invoice;

    Invoices.actions.complete(locId);
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
    Invoices.actions.selectInvoiceId(null);
  },

  makeDraft(e) {
    e.stopPropagation();
    const locId = this.state.locationId;
    const invoice = this.state.invoice;

    Invoices.actions.makeDraft(locId);
  },

  render() {
    const invoice = this.state.invoice;
    const invoiceStatuses = this.state.invoiceStatuses;

    if (!invoice) {
      return <LoaderLocal/>;
    }

    const amount = (Math.round(invoice.getIn(['Sums', 'All', 'PriceInclVAT']) * 100)
                    / 100).toFixed(2);
    const name = invoice.getIn(['User', 'FirstName']) +
                 ' ' + invoice.getIn(['User', 'LastName']);
    const timeFrame = '' + invoice.get('Month') + '/' +
                           invoice.get('Year');

    return (
      <div id="inv-header">
        <div className="row">
          <div className="col-xs-3">
            <h3>{name}</h3>
          </div>
          <div className="col-xs-3 inv-time-frame">
            <h3>Invoice {timeFrame}</h3>
          </div>
          <div className="col-xs-3 inv-amount">
            <h3>Amount {amount} €</h3>
          </div>
          <div className="col-xs-3">
            <button type="button"
                    title="Close"
                    onClick={this.hide}>
              <i className="fa fa-close"/>
            </button>
            {this.state.invoiceActions.get('PushDraft') ?
              <button type="button"
                      onClick={this.makeDraft}
                      title="Make Draft">
                <i className="fa fa-refresh"/>
              </button> : null
            }
            {this.state.invoiceActions.get('Save') ?
              <button type="button"
                      onClick={this.save}
                      title="Save">
                <i className="fa fa-check"/>
              </button> : null
            }
            {this.state.invoiceActions.get('Freeze') ?
              <button type="button"
                      onClick={this.complete}
                      title="Freeze">
                <i className="fa fa-cart-arrow-down"/>
              </button> : null
            }
            {this.state.invoiceActions.get('Send') ?
              <button type="button"
                      onClick={this.send}
                      title="Send">
                <i className="fa fa-send"/>
              </button> : null
            }
            {this.state.invoiceActions.get('SendCanceled') ?
              <button type="button"
                      onClick={this.sendCanceled}
                      title="Send Canceled">
                <i className="fa fa-send"/>
              </button> : null
            }
            {this.state.invoiceActions.get('Cancel') ?
              <button type="button"
                      onClick={this.cancel}
                      title="Cancel">
                <i className="fa fa-ban"/>
              </button> : null
            }
          </div>
        </div>
        <div className="row">
          <div className="col-xs-3">
            Fastbill No: {invoice.getIn(['User', 'ClientId'])}
          </div>
          <div className="col-xs-3">
            {invoice.get('FastbillNo') ?
              'Invoice No: ' + invoice.get('FastbillNo') :
              'Draft'
            }
          </div>
          <div className="col-xs-3">
            Incl. {this.state.vatPercent}% VAT
          </div>
          <div className="col-xs-3">
          </div>
        </div>
      </div>
    );
  },

  save(e) {
    e.stopPropagation();
    const locId = this.state.locationId;
    const userId = this.state.uid;
    const month = this.state.MonthlySums
                      .get('selected').get('month');
    const year = this.state.MonthlySums
                      .get('selected').get('year');

    Invoices.actions.save(locId, {month, year, userId});
  },

  send(e) {
    e.stopPropagation();

    Invoices.actions.send();
  },

  sendCanceled(e) {
    e.stopPropagation();

    Invoices.actions.sendCanceled();
  }
});


var Invoice = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      editPurchaseId: Invoices.getters.getEditPurchaseId,
      invoice: Invoices.getters.getInvoice,
      invoiceActions: Invoices.getters.getInvoiceActions,
      invoiceStatuses: Invoices.getters.getInvoiceStatuses,
      location: LocationGetters.getLocation,
      locationId: LocationGetters.getLocationId,
      MonthlySums: Invoices.getters.getMonthlySums,
      uid: getters.getUid,
      userMemberships: Invoices.getters.getUserMemberships
    };
  },

  render() {
    const invoice = this.state.invoice;
    const invoiceStatuses = this.state.invoiceStatuses;
    console.log('invoiceActions=', this.state.invoiceActions.toJS());
    if (invoice) {


      return (
        <div className="inv-invoice"
             onClick={this.stopPropagation}>
          <Header/>
          <div id="inv-body">
            <BillTable bill={invoice}/>
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

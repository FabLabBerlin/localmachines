var BillTable = require('../../UserProfile/BillTable');
var getters = require('../../../getters');
var Invoices = require('../../../modules/Invoices');
var LoaderLocal = require('../../LoaderLocal');
var LocationGetters = require('../../../modules/Location/getters');
var Membership = require('../../UserProfile/Membership');
var moment = require('moment');
var React = require('react');
var reactor = require('../../../reactor');
var toastr = require('../../../toastr');


var Invoice = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      invoice: Invoices.getters.getInvoice,
      invoiceStatuses: Invoices.getters.getInvoiceStatuses,
      location: LocationGetters.getLocation,
      locationId: LocationGetters.getLocationId,
      MonthlySums: Invoices.getters.getMonthlySums,
      uid: getters.getUid,
      userMemberships: Invoices.getters.getUserMemberships
    };
  },

  hide() {
    Invoices.actions.selectUserId(null);
  },

  makeDraft(e) {
    e.stopPropagation();
    const locId = this.state.locationId;
    const userId = this.state.uid;
    const month = this.state.MonthlySums
                      .get('selected').get('month');
    const year = this.state.MonthlySums
                      .get('selected').get('year');

    toastr.info('Invoice#makeDraft()');
    Invoices.actions.makeDraft(locId, {month, year, userId});
  },

  render() {
    const invoice = this.state.invoice;
    const invoiceStatuses = this.state.invoiceStatuses;

    console.log('invoiceStatuses:', invoiceStatuses);

    if (invoice && this.state.userMemberships) {
      const amount = (Math.round(invoice.Sums.All.PriceInclVAT * 100)
                      / 100).toFixed(2);
      const name = invoice.User.FirstName + ' ' + invoice.User.LastName;
      const timeFrame = '' + invoice.Interval.MonthFrom + '/' +
                             invoice.Interval.YearFrom;
      console.log('invoice=', invoice);

      return (
        <div className="inv-invoice-container" onClick={this.hide}>
          <div className="inv-invoice-background"/>
          <div className="inv-invoice-aligner">
            <div className="inv-invoice">
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
                          onClick={this.makeDraft}
                          title="Make Draft">
                    <i className="fa fa-pencil"/>
                  </button>
                  <button type="button"
                          title="Send">
                    <i className="fa fa-send"/>
                  </button>
                </div>
              </div>
              <h5>Activations</h5>
              <BillTable bill={invoice}/>
              <h5>Memberships</h5>
              <Membership memberships={this.state.userMemberships.Data}/>
            </div>
          </div>
        </div>
      );
    } else {
      return <LoaderLocal/>;
    }
  }

});

export default Invoice;

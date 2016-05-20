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

    if (invoice && this.state.userMemberships) {
      const name = invoice.User.FirstName + ' ' + invoice.User.LastName;

      return (
        <div className="inv-invoice-container" onClick={this.hide}>
          <div className="inv-invoice-background"/>
          <div className="inv-invoice-aligner">
            <div className="inv-invoice">
              User Invoice for {name}
              <div>
                <button type="button" onClick={this.makeDraft}>
                  <i className="fa fa-pencil"/> Make Draft
                </button>
                <button type="button">
                  <i className="fa fa-send"/> Send
                </button>
              </div>
              <Membership memberships={this.state.userMemberships.Data}/>
              <BillTable bill={invoice}/>
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

var BillTable = require('../../UserProfile/BillTable');
var Invoices = require('../../../modules/Invoices');
var LoaderLocal = require('../../LoaderLocal');
var LocationGetters = require('../../../modules/Location/getters');
var Membership = require('../../UserProfile/Membership');
var moment = require('moment');
var React = require('react');
var reactor = require('../../../reactor');


var Invoice = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      invoice: Invoices.getters.getInvoice,
      location: LocationGetters.getLocation,
      locationId: LocationGetters.getLocationId,
      MonthlySums: Invoices.getters.getMonthlySums,
      userMemberships: Invoices.getters.getUserMemberships
    };
  },

  hide() {
    Invoices.actions.selectUserId(null);
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

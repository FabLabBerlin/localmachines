var BillTable = require('../../UserProfile/BillTable');
var Invoices = require('../../../modules/Invoices');
var LoaderLocal = require('../../LoaderLocal');
var LocationGetters = require('../../../modules/Location/getters');
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
      MonthlySums: Invoices.getters.getMonthlySums
    };
  },

  hide() {
    Invoices.actions.selectUserId(null);
  },

  render() {
    if (this.state.invoice) {
      console.log('this.state.invoice:', this.state.invoice);
      return (
        <div className="inv-invoice-container" onClick={this.hide}>
          <div className="inv-invoice-background"/>
          <div className="inv-invoice-aligner">
            <div className="inv-invoice">
              User Invoice
              <BillTable bill={this.state.invoice}/>
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

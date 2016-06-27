var Invoice = require('./Invoice');
var Invoices = require('../../../modules/Invoices');
var LoaderLocal = require('../../LoaderLocal');
var LocationGetters = require('../../../modules/Location/getters');
var moment = require('moment');
var React = require('react');
var reactor = require('../../../reactor');


var List = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      locationId: LocationGetters.getLocationId,
      MonthlySums: Invoices.getters.getMonthlySums
    };
  },

  render() {
    const summaries = this.props.summaries.sortBy((sum) => {
      return (sum.getIn(['User', 'FirstName'])
           + ' ' + sum.getIn(['User', 'LastName']))
           .toLowerCase();
    });
    const total = (Math.round(summaries.reduce((result, monthlySum) => {
      return result + monthlySum.get('Amount');
    }, 0) * 100) / 100).toFixed(2);
    console.log('total=', total);

    return (
      <div>
        {this.state.MonthlySums.getIn(['selected', 'invoiceId']) ? (
          <Invoice/>
        ) : null}
        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th>
                <label><input type="checkbox"/> All</label>
              </th>
              <th className="text-center">
                <label>No.</label>
              </th>
              <th className="text-center">
                <label>Canceled</label>
              </th>
              <th className="text-center">
                <label>Users</label>
              </th>
              <th className="text-center">
                <label>Paid</label>
              </th>
              <th className="text-center">
                <label>Total</label>
              </th>
            </tr>
          </thead>
          <tbody>
            {summaries.map((sum, i) => {
              const name = sum.getIn(['User', 'FirstName'])
                         + ' ' + sum.getIn(['User', 'LastName']);
              const amount = (Math.round(sum.get('Total') * 100)
                                 / 100)
                                 .toFixed(2);

              return (
                <tr key={i} onClick={this.select.bind(this, sum.get('User'), sum.get('Id'))}>
                  <td><input type="checkbox"/></td>
                  <td className="text-right">
                    {sum.get('FastbillNo') || 'Draft'}
                  </td>
                  <td className="text-center">
                    {sum.get('Canceled') ?
                      <i className="fa fa-check"/> :
                      '-'
                    }
                  </td>
                  <td>{name}</td>
                  <td className="text-center">
                    {(moment(sum.get('PaidDate')).unix() > 0) ?
                      <i className="fa fa-check"/> :
                      '-'
                    }
                  </td>
                  <td className="text-right">{amount} €</td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    );
  },

  select(user, invoiceId) {
    const month = this.state.MonthlySums
                      .get('selected').get('month');
    const year = this.state.MonthlySums
                      .get('selected').get('year');
    const userId = user.get('Id');

    Invoices.actions.fetchInvoice(this.state.locationId, {
      invoiceId: invoiceId
    });
    Invoices.actions.selectInvoiceId(invoiceId);
  }

});

export default List;

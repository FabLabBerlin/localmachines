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
      checkedAll: Invoices.getters.getCheckedAll,
      locationId: LocationGetters.getLocationId,
      MonthlySums: Invoices.getters.getMonthlySums,
      checkStatus: Invoices.getters.getCheckStatus
    };
  },

  check(invoiceId, e) {
    console.log('e=', e);
    e.stopPropagation();
    Invoices.actions.check(invoiceId);
  },

  checkAll(e) {
    e.stopPropagation();
    Invoices.actions.checkAll();
  },

  render() {
    const summaries = this.props.summaries.sortBy((sum) => {
      return (sum.getIn(['User', 'FirstName'])
           + ' ' + sum.getIn(['User', 'LastName']))
           .toLowerCase();
    });

    if (this.state.MonthlySums.getIn(['selected', 'invoiceId'])) {
      return <Invoice/>;
    }

    return (
      <div>
        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th>
                <input type="checkbox"
                       checked={this.state.checkedAll}
                       onChange={this.checkAll}/>
              </th>
              <th>
                <select onChange={this.selectStatus}
                        value={this.state.checkStatus}>
                  <option value="all">All</option>
                  <option value="draft">Draft</option>
                  <option value="outgoing">Outgoing</option>
                </select>
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
              const click = this.select.bind(this, sum.get('User'), sum.get('Id'));

              return (
                <tr key={i}>
                  <td>
                    <input type="checkbox"
                           checked={sum.get('checked')}
                           onChange={this.check.bind(this, sum.get('Id'))}/>
                  </td>
                  <td onClick={click}/>
                  <td className="text-right" onClick={click}>
                    {sum.get('FastbillNo') || 'Draft'}
                  </td>
                  <td className="text-center" onClick={click}>
                    {sum.get('Canceled') ?
                      <i className="fa fa-check"/> :
                      '-'
                    }
                  </td>
                  <td onClick={click}>{name}</td>
                  <td className="text-center" onClick={click}>
                    {(moment(sum.get('PaidDate')).unix() > 0) ?
                      <i className="fa fa-check"/> :
                      '-'
                    }
                  </td>
                  <td className="text-right" onClick={click}>{amount} €</td>
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
  },

  selectStatus(e) {
    console.log('e.target.value=', e.target.value);
    Invoices.actions.checkSetStatus(e.target.value);
  }

});

export default List;

var _ = require('lodash');
var Invoice = require('./Invoice');
var Invoices = require('../../../modules/Invoices');
var LoaderLocal = require('../../LoaderLocal');
var LocationGetters = require('../../../modules/Location/getters');
var moment = require('moment');
var React = require('react');
var reactor = require('../../../reactor');
var Settings = require('../../../modules/Settings');
var util = require('./util');


var HeadColumn = React.createClass({

  asc() {
    var asc;
    const sorting = this.props.sorting;

    if (sorting && sorting.get('column') === this.props.attribute) {
      asc = sorting.get('asc');
    }

    return asc;
  },

  onClick() {
    var asc = this.asc();

    if (asc) {
      asc = false;
    } else if (asc === false) {
      asc = undefined;
    } else {
      asc = true;
    }

    this.props.toggle(this.props.attribute, asc);
  },

  render() {
    return (
      <th className="text-center">
        <label onClick={this.onClick}>
          {this.props.label} <SortIndicator asc={this.asc()}/>
        </label>
      </th>
    );
  }
});


var SortIndicator = React.createClass({
  render() {
    if (this.props.asc) {
      return <i className="fa fa-sort-asc"/>;
    } else if (_.isUndefined(this.props.asc)) {
      return <i className="fa fa-sort"/>;
    } else {
      return <i className="fa fa-sort-desc"/>;
    }
  }
});


var List = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      checkedAll: Invoices.getters.getCheckedAll,
      currency: Settings.getters.getCurrency,
      locationId: LocationGetters.getLocationId,
      MonthlySums: Invoices.getters.getMonthlySums,
      checkStatus: Invoices.getters.getCheckStatus
    };
  },

  check(invoiceId, e) {
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

    var sorted;

    const sorting = this.state.MonthlySums.get('sorting');

    if (sorting) {
      sorted = summaries.sortBy((inv) => {
        const c = sorting.get('column').split('.');
        const v = inv.getIn(c);

        if (c === 'FastbillNo') {
          return parseInt(v, 10) || -1;
        } else {
          return v;
        }
      });
      if (sorting.get('asc') === false) {
        sorted = sorted.reverse();
      }
    } else {
      sorted = summaries;
    }

    return (
      <div>
        {this.state.MonthlySums.getIn(['selected', 'invoiceId']) ?
          <Invoice/> : null}
        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th>
                <input checked={this.state.checkedAll}
                       className="invs-select-all"
                       onChange={this.checkAll}
                       type="checkbox"/>
              </th>
              <th>
                <select className="invs-select-all"
                        onChange={this.selectStatus}
                        value={this.state.checkStatus}>
                  <option value="all">All</option>
                  <option value="draft">Draft</option>
                  <option value="outgoing">Outgoing</option>
                </select>
              </th>
              <HeadColumn label="No."
                          attribute="FastbillNo"
                          sorting={sorting}
                          toggle={this.toggleSorting}/>
              <HeadColumn label="Status"
                          attribute="Status"
                          sorting={sorting}
                          toggle={this.toggleSorting}/>
              <HeadColumn label="Canceled"
                          attribute="Canceled"
                          sorting={sorting}
                          toggle={this.toggleSorting}/>
              <HeadColumn label="User"
                          attribute="Name"
                          sorting={sorting}
                          toggle={this.toggleSorting}/>
              <HeadColumn label="Paid"
                          attribute="PaidDate"
                          sorting={sorting}
                          toggle={this.toggleSorting}/>
              <HeadColumn label="No Auto Invoicing"
                          attribute="User.NoAutoInvoicing"
                          sorting={sorting}
                          toggle={this.toggleSorting}/>
              <HeadColumn label="Total"
                          attribute="Total"
                          sorting={sorting}
                          toggle={this.toggleSorting}/>
            </tr>
          </thead>
          <tbody>
            {sorted.map((sum, i) => {
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
                    {sum.get('FastbillNo') || '-'}
                  </td>
                  <td className="text-right" onClick={click}>
                    {util.statusInfo(sum)}
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
                  <td className="text-center">
                    {sum.getIn(['User', 'NoAutoInvoicing']) ?
                      <i className="fa fa-check"/> :
                      '-'
                    }
                  </td>
                  <td className="text-right" onClick={click}>{amount} {this.state.currency}</td>
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
    Invoices.actions.checkSetStatus(e.target.value);
  },

  toggleSorting(colName, asc) {
    Invoices.actions.sortBy(colName, asc);
  }

});

export default List;

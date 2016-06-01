var Invoice = require('./Invoice');
var Invoices = require('../../../modules/Invoices');
var LoaderLocal = require('../../LoaderLocal');
var LocationGetters = require('../../../modules/Location/getters');
var moment = require('moment');
var React = require('react');
var reactor = require('../../../reactor');


var Summaries = React.createClass({

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
        {this.state.MonthlySums.getIn(['selected', 'userId']) ? (
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
                <tr key={i} onClick={this.select.bind(this, sum.get('User'))}>
                  <td><input type="checkbox"/></td>
                  <td className="text-right">
                    {sum.get('FastbillNo') || 'Draft'}
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

  select(user) {
    const month = this.state.MonthlySums
                      .get('selected').get('month');
    const year = this.state.MonthlySums
                      .get('selected').get('year');
    const userId = user.get('Id');

    Invoices.actions.fetchFastbillStatuses(this.state.locationId, {
      year: year,
      month: month,
      userId: userId
    });
    Invoices.actions.fetchUser(this.state.locationId, {
      year: year,
      month: month,
      userId: userId
    });
    Invoices.actions.fetchUserMemberships(this.state.locationId, {
      userId: userId
    });
    Invoices.actions.selectUserId(userId);
  }

});

export default Summaries;

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
        <table>
          <thead>
            <tr>
              <th>User</th>
              <th>Invoiced?</th>
              <th>Paid?</th>
              <th>Amount (EUR)</th>
            </tr>
          </thead>
          <tbody>
            {summaries.map((sum, i) => {
              if (sum.getIn(['User', 'Id']) !== 19) {
                return undefined;
              }
              const name = sum.getIn(['User', 'FirstName'])
                         + ' ' + sum.getIn(['User', 'LastName']);
              const amount = (Math.round(sum.get('Amount') * 100)
                                 / 100)
                                 .toFixed(2);

              return (
                <tr key={i} onClick={this.select.bind(this, sum.get('User'))}>
                  <td>{name}</td>
                  <td></td>
                  <td></td>
                  <td className="text-right">{amount}</td>
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

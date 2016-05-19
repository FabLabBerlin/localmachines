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

    return (
      <div>
        {this.state.MonthlySums.getIn(['selected', 'userId']) ? (
          <Invoice/>
        ) : null}
        <table>
          <thead>
            <tr>
              <th>User</th>
              <th>Amount (EUR)</th>
              <th>Invoiced?</th>
              <th>Paid?</th>
            </tr>
          </thead>
          <tbody>
            {summaries.map((sum, i) => {
              const name = sum.getIn(['User', 'FirstName'])
                         + ' ' + sum.getIn(['User', 'LastName']);
              const amount = (Math.round(sum.get('Amount') * 100)
                                 / 100)
                                 .toFixed(2);

              return (
                <tr key={i} onClick={this.select.bind(this, sum.get('User'))}>
                  <td>{name}</td>
                  <td className="text-right">{amount}</td>
                  <td></td>
                  <td></td>
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

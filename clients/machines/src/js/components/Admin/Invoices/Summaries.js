var Invoices = require('../../../modules/Invoices');
var LoaderLocal = require('../../LoaderLocal');
var LocationGetters = require('../../../modules/Location/getters');
var moment = require('moment');
var React = require('react');
var reactor = require('../../../reactor');


var Summaries = React.createClass({

  render() {
    const summaries = this.props.summaries.sortBy((sum) => {
      return (sum.getIn(['User', 'FirstName'])
           + ' ' + sum.getIn(['User', 'LastName']))
           .toLowerCase();
    });

    return (
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
              <tr key={i}>
                <td>{name}</td>
                <td className="text-right">{amount}</td>
                <td></td>
                <td></td>
              </tr>
            );
          })}
        </tbody>
      </table>
    );
  }

});

export default Summaries;

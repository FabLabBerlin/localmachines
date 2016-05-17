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


var Month = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      location: LocationGetters.getLocation,
      locationId: LocationGetters.getLocationId,
      monthlySummaries: Invoices.getters.getMonthlySummaries
    };
  },

  componentWillMount() {
    console.log('Month#componentWillMount');

    const t = this.props.t;
    const month = this.state.monthlySummaries
                      .get('selected').get('month');
    const year = this.state.monthlySummaries
                      .get('selected').get('year');
    const selected = month === t.month() + 1 &&
                     year === t.year();
    const summaries = this.state.monthlySummaries.getIn([t.year(), t.month() + 1]);

    if (selected && !summaries) {
      Invoices.actions.fetchMonthlySummaries(this.state.locationId, {
        month: t.month() + 1,
        year: t.year()
      });
    }
  },

  render() {
    const t = this.props.t;
    const month = this.state.monthlySummaries
                      .get('selected').get('month');
    const year = this.state.monthlySummaries
                      .get('selected').get('year');
    const selected = month === t.month() + 1 &&
                     year === t.year();
    const summaries = this.state.monthlySummaries.getIn([year, month]);

    if (selected && !summaries) {
      return <LoaderLocal/>;
    } else {
      return (
        <div>
          <h3 onClick={this.select}>{t.format('MMMM YYYY')}</h3>
          {
            selected ? (
              <Summaries summaries={summaries}/>
            ) : null
          }
        </div>
      );
    }
  },

  select() {
    Invoices.actions.setSelectedMonth({
      month: this.props.t.month() + 1,
      year: this.props.t.year()
    });
  }

});


var InvoicesView = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      location: LocationGetters.getLocation,
      locationId: LocationGetters.getLocationId,
      monthlySummaries: Invoices.getters.getMonthlySummaries
    };
  },

  render() {
    var t = moment();
    var nodes = [];

    for (var i = 0; i < 12; i++) {
      t = t.clone().subtract(1, 'months');
      nodes.push(
        <Month key={i} t={t}/>
      );
    }
    return (
      <div className="container-fluid">
        <h2>Invoices</h2>
        {nodes}
      </div>
    );
  }

});

export default InvoicesView;

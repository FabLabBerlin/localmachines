var Invoices = require('../../../modules/Invoices');
var LoaderLocal = require('../../LoaderLocal');
var LocationGetters = require('../../../modules/Location/getters');
var moment = require('moment');
var React = require('react');
var reactor = require('../../../reactor');
var SettingsActions = require('../../../modules/Settings/actions');
var Summaries = require('./Summaries');


var Month = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      location: LocationGetters.getLocation,
      locationId: LocationGetters.getLocationId,
      MonthlySums: Invoices.getters.getMonthlySums
    };
  },

  render() {
    const t = this.props.t;
    const month = this.state.MonthlySums
                      .get('selected').get('month');
    const year = this.state.MonthlySums
                      .get('selected').get('year');
    const selected = month === t.month() + 1 &&
                     year === t.year();
    const summaries = this.state.MonthlySums.getIn([year, month]);

    return (
      <div className="inv-monthly-sums">
        <h3 onClick={this.select}>{t.format('MMMM YYYY')}</h3>
        {
          selected ? (
            summaries ? (
              <Summaries summaries={summaries}/>
            ) : (
              <LoaderLocal/>
            )
          ) : null
        }
      </div>
    );
  },

  select() {
    const t = this.props.t;
    const summaries = this.state.MonthlySums.getIn([t.year(), t.month() + 1]);

    Invoices.actions.setSelectedMonth({
      month: t.month() + 1,
      year: t.year()
    });

    if (!summaries) {
      Invoices.actions.fetchMonthlySums(this.state.locationId, {
        month: t.month() + 1,
        year: t.year()
      });
    }
  }

});


var InvoicesView = React.createClass({

  mixins: [ reactor.ReactMixin ],

  componentWillMount() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    var t = moment();
    Invoices.actions.fetchMonthlySums(this.state.locationId, {
      month: t.month(),
      year: t.year()
    });
    SettingsActions.loadSettings({locationId});
  },

  getDataBindings() {
    return {
      location: LocationGetters.getLocation,
      locationId: LocationGetters.getLocationId,
      MonthlySums: Invoices.getters.getMonthlySums
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
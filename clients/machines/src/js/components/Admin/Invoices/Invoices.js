var getters = require('../../../getters');
var Invoices = require('../../../modules/Invoices');
var List = require('./List');
var LoaderLocal = require('../../LoaderLocal');
var LocationActions = require('../../../actions/LocationActions');
var LocationGetters = require('../../../modules/Location/getters');
var moment = require('moment');
var React = require('react');
var reactor = require('../../../reactor');
var SettingsActions = require('../../../modules/Settings/actions');


var Month = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      location: LocationGetters.getLocation,
      locationId: LocationGetters.getLocationId,
      MonthlySums: Invoices.getters.getMonthlySums,
      uid: getters.getUid
    };
  },

  componentWillMount() {
    LocationActions.loadUserLocations(this.state.uid);
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

    if (!summaries && selected) {
      return <LoaderLocal/>;
    }

    var total;
    if (selected) {
      total = (Math.round(this.state.MonthlySums.getIn([
        t.year(),
        t.month() + 1
      ]).reduce((result, monthlySum) => {
        return result + monthlySum.get('Total');
      }, 0) * 100) / 100).toFixed(2);
      console.log('total=', total);
    }

    return (
      <div className={'inv-monthly-sums ' + (selected ? 'selected' : '')}>
        <div className="row">
          <div className="col-md-6">
            <h3 onClick={this.select}>{t.format('MMMM YYYY')}</h3>
          </div>
          <div className="col-md-6 text-right">
            {selected ? ('Sum total:' + total + ' €') : null}
          </div>
        </div>
        {
          selected ? (
            summaries ? (
              <List summaries={summaries}/>
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
      month: t.month() + 1,
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
      nodes.push(
        <Month key={i} t={t}/>
      );
      t = t.clone().subtract(1, 'months');
    }
    return (
      <div className="container-fluid">
        <h2>Invoices EXPERIMENTAL</h2>
        {nodes}
      </div>
    );
  }

});

export default InvoicesView;

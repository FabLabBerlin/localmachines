var Invoices = require('../../../modules/Invoices');
var LocationGetters = require('../../../modules/Location/getters');
var moment = require('moment');
var React = require('react');
var reactor = require('../../../reactor');


var InvoicesView = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      location: LocationGetters.getLocation,
      locationId: LocationGetters.getLocationId,
      monthlySummaries: Invoices.getters.getMonthlySummaries
    };
  },

  componentWillMount() {
    const t = moment();
    Invoices.actions.fetchMonthlySummaries(this.state.locationId, {
      month: t.month() + 1,
      year: t.year()
    });
  },

  render() {
    var t = moment();
    var nodes = [];

    for (var i = 0; i < 12; i++) {
      t = t.clone().subtract(1, 'months');
      nodes.push(
        <div key={i}>
          <h3>{t.format('MMMM YYYY')}</h3>
        </div>
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

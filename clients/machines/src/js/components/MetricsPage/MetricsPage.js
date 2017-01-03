var Categories = require('../../modules/Categories');
var Invoices = require('../../modules/Invoices');
var Location = require('../../modules/Location');
var Machines = require('../../modules/Machines');
var getters = require('../../../getters');
var React = require('react');
var reactor = require('../../../reactor');


var MetricsPage = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      activations: Machines.getters.getActivations,
      categories: Categories.getters.getAll,
      invoices: Invoices.getters.getAllInvoices,
      locationId: Location.getters.getLocationId,
      machines: Machines.getters.getMachines,
      machinesById: Machines.getters.getMachinesById,
      myMachines: Machines.getters.getMyMachines,
      upcomingReservations: getters.getUpcomingReservationsByMachineId
    };
  },

  componentWillMount() {
    const locId = reactor.evaluateToJS(Location.getters.getLocationId);

    Invoices.actions.fetchInvoices(locId);
  },

  render() {
    return (
      <div>
        
      </div>
    );
  }
});

export default MetricsPage;

var React = require('react');


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
    Invoices.actions.fetchInvoices
  },

  render() {
    return (
      <div>
        
      </div>
    );
  }
});

export default MetricsPage;

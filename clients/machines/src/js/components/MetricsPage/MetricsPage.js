import Categories from '../../modules/Categories';
import Invoices from '../../modules/Invoices';
import Location from '../../modules/Location';
import Machines from '../../modules/Machines';
import getters from '../../../getters';
import React from 'react';
import reactor from '../../../reactor';


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

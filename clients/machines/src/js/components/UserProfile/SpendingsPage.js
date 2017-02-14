import BillTables from './BillTables';
import getters from '../../getters';
import LoaderLocal from '../LoaderLocal';
import Location from '../../modules/Location';
import LoginActions from '../../actions/LoginActions';
import MachineActions from '../../actions/MachineActions';
import Machines from '../../modules/Machines';
import React from 'react';
import reactor from '../../reactor';
import SettingsActions from '../../modules/Settings/actions';
import UserActions from '../../actions/UserActions';


var SpendingsPage = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      user: getters.getUser,
      machines: Machines.getters.getMachines,
      bill: getters.getBill,
      memberships: getters.getMemberships
    };
  },

  componentDidMount() {
    const locationId = reactor.evaluateToJS(Location.getters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);
    MachineActions.apiGetUserMachines(locationId, uid);
    UserActions.fetchUser(uid);
    UserActions.fetchBill(locationId, uid);
    Location.actions.loadUserLocations(uid);
    SettingsActions.loadSettings({locationId});
  },

  render() {
    if (this.state.bill) {
      return (
        <div className="container">
          <BillTables bill={this.state.bill}
                      membership={this.state.membership}/>
        </div>
      );
    } else {
      return <LoaderLocal/>;
    }
  }
});

export default SpendingsPage;

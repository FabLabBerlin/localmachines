var BillTables = require('./BillTables');
var getters = require('../../getters');
var LoaderLocal = require('../LoaderLocal');
var Location = require('../../modules/Location');
var LoginActions = require('../../actions/LoginActions');
var MachineActions = require('../../actions/MachineActions');
var Machines = require('../../modules/Machines');
var React = require('react');
var reactor = require('../../reactor');
var SettingsActions = require('../../modules/Settings/actions');
var UserActions = require('../../actions/UserActions');


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

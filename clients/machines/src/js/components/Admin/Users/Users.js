var _ = require('lodash');
var getters = require('../../../getters');
var LocationGetters = require('../../../modules/Location/getters');
var MachineActions = require('../../../actions/MachineActions');
var Machines = require('../../../modules/Machines');
var Navigation = require('react-router').Navigation;
var React = require('react');
var reactor = require('../../../reactor');


var UsersView = React.createClass({

  mixins: [ Navigation, reactor.ReactMixin ],

  getDataBindings() {
    return {
      location: LocationGetters.getLocation,
      machines: Machines.getters.getMachines
    };
  },

  componentWillMount() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);
    MachineActions.apiGetUserMachines(locationId, uid);
  },

  render() {
    var users = _.sortBy(this.state.users.toJS(), m => m.Name);
    return (
      <div className="container-fluid">
        <div className="row">
          <div className="col-xs-1">
          </div>
          <div className="col-xs-11">
          </div>
        </div>
      </div>
    );
  }

});

export default UsersView;

var getters = require('../../../getters');
var LoaderLocal = require('../../LoaderLocal');
var Machines = require('../../../modules/Machines');
var React = require('react');
var reactor = require('../../../reactor');


var OccupiedBy = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      machineUsers: Machines.getters.getMachineUsers,
      uid: getters.getUid
    };
  },

  render() {
    if (this.props.activation) {
      const uid = this.props.activation.get('UserId');
      if (uid === this.state.uid) {
        return <div/>;
      }

      if (this.state.machineUsers) {
        var users = this.state.machineUsers;
        var user = users.get(uid) || {};

        return (
          <div className="m-indicator m-occupied-by">
            <div>Occupied by</div>
            <div>{user.FirstName} {user.LastName}</div>
          </div>
        );
      } else {
        return <LoaderLocal/>;
      }
    } else {
      return <div/>;
    }
  }
});

export default OccupiedBy;

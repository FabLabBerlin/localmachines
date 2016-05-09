var _ = require('lodash');
var getters = require('../../getters');
var LoaderLocal = require('../LoaderLocal');
var LocationGetters = require('../../modules/Location/getters');
var Machine = require('./Machine/Machine');
var React = require('react');
var reactor = require('../../reactor');


var MachineList = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      location: LocationGetters.getLocation,
      machines: getters.getMachines
    };
  },

  render() {
    let activation = this.props.activation;
    var MachineNode;
    if(this.state.location && this.props.machines && _.size(this.props.machines) > 0) {
      var machines = _.filter(this.props.machines, function(machine) {
        return machine.Visible && !machine.Archived;
      });
      MachineNode = _.map(machines, function(machine) {
        if (machine.LocationId === this.state.location.Id) {
          let activationProps = false;
          let isMachineBusy = false;
          let isSameUser = false;
          if (activation) {
            activation.forEach(function(a, i) {
              if( machine.Id === a.get('MachineId') ) {
                isMachineBusy = true;
                isSameUser = this.props.user.get('Id') === a.get('UserId');
                activationProps = a.toJS();
                return false;
              }
            }.bind(this));
          }
          return (
            <Machine
              key={machine.Id}
              machine={machine}
              user={this.props.user}
              busy={isMachineBusy}
              sameUser={isSameUser}
              activation={activationProps}
            />
          );
        }
      }.bind(this));
    } else {
      return <LoaderLocal/>;
    }
    return (
      <div className="machines">
        <div className="container-fluid">
          <h2>Available Machines</h2>
        </div>
        {MachineNode}
      </div>
    );
  }
});

export default MachineList;

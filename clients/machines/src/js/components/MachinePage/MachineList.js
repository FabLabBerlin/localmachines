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
    console.log('activations:', activation);
    var MachineNode;
    if (this.state.location && this.props.machines) {
      var machines = _.filter(this.props.machines, function(machine) {
        return machine.Visible && !machine.Archived;
      });
      MachineNode = _.map(machines, function(machine) {
        if (machine.LocationId === this.state.location.Id) {
          let activationProps;
          if (activation) {
            activation.forEach(function(a, i) {
              if( machine.Id === a.get('MachineId') ) {
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
              activation={activationProps}
            />
          );
        }
      }.bind(this));
    } else {
      console.log('bla');
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

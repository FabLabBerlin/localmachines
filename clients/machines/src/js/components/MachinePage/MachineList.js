var _ = require('lodash');
var getters = require('../../getters');
var LoaderLocal = require('../LoaderLocal');
var LocationGetters = require('../../modules/Location/getters');
var Machine = require('./Machine/Machine');
var Machines = require('../../modules/Machines');
var React = require('react');
var reactor = require('../../reactor');


var MachineList = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      locationId: LocationGetters.getLocationId,
      machines: Machines.getters.getMachines
    };
  },

  render() {
    let activation = this.props.activation;
    var MachineNode;
    if (this.state.locationId && this.props.machines) {
      var machines = this.props.machines.toList().filter((machine) => {
        return machine.get('Visible') && !machine.get('Archived');
      });

      machines = _.sortBy(machines.toJS(), m => {
        const typeName = {
          0: 'Other',
          1: '3D Printer',
          2: 'CNC mill',
          3: 'Heatpress',
          4: 'Knitting Machine',
          5: 'Lasercutters',
          6: 'Vinylcutter'
        }[m.TypeId];

        return (typeName + m.Name).toLowerCase();
      });
      MachineNode = _.map(machines, function(machine) {
        if (machine.LocationId === this.state.locationId) {
          let activationProps;
          if (activation) {
            activation.forEach(function(a, i) {
              if( machine.Id === a.get('MachineId') ) {
                activationProps = a.toJS();
                return false;
              }
            }.bind(this));
          }
          if (!machine.Locked) {
            return (
              <Machine
                key={machine.Id}
                machine={machine}
                user={this.props.user}
                activation={activationProps}
              />
            );
          }
        }
      }.bind(this));
    } else {
      return <LoaderLocal/>;
    }
    return (
      <div className="machines">
        {MachineNode}
      </div>
    );
  }
});

export default MachineList;

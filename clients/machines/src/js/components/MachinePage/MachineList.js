var _ = require('lodash');
var getters = require('../../getters');
var React = require('react');
var reactor = require('../../reactor');
var Machine = require('./Machine/Machine');


var MachineList = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      location: getters.getLocation,
      machines: getters.getMachines
    };
  },

  render() {
    let activation = this.props.activation;
    var MachineNode;
    if(this.state.location && this.props.machines && _.size(this.props.machines) > 0) {
      var machines = _.filter(this.props.machines, function(machine) {
        return machine.Visible;
      });
      console.log('this.state.location.Id:', this.state.location.Id);
      MachineNode = _.map(machines, function(machine) {
        if (machine.LocationId === this.state.location.Id) {
          let activationProps = false;
          let isMachineBusy = false;
          let isSameUser = false;
          if (activation) {
            _.each(activation, function(a, i) {
              if( machine.Id === activation[i].MachineId ) {
                isMachineBusy = true;
                isSameUser = this.props.user.get('Id') === activation[i].UserId;
                activationProps = activation[i];
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
      return (
        <div className="loader-local">
          <div className="spinner">
            <i className="fa fa-cog fa-spin"></i>
          </div>
        </div>
      );
    }
    return (
      <div className="machines">
        <div className="container-fluid">
          <h2>Available machines</h2>
        </div>
        {MachineNode}
      </div>
    );
  }
});

export default MachineList;

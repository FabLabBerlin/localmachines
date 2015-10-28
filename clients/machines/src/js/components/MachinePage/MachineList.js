var _ = require('lodash');
var getters = require('../../getters');
var React = require('react');
var reactor = require('../../reactor');
var MachineChooser = require('./MachineChooser');


/*
 * MachineList component:
 * Create a list of MachineChooser
 * Prepare the props for the MachineChooser to display the right MachineDiv
 */
var MachineList = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      machineInfo: getters.getMachineInfo
    };
  },

  /*
   * Render a list of MachineChooser
   * Prepare props for each MachineChooser
   * Parse the activation array and the machine array
   *    and distribute the information concerning a machine
   *    in only a child component
   * @activationProps: false is no activation (ie FreeMachine)
   *    activation object of the machine otherwise
   * @isMachineBusy: true if busy, false otherwise
   * @sameUser: true if the user of the machine and the one on the activation are the same
   *    false otherwise
   */
  render() {
    let activation = this.props.activation;
    var MachineNode;
    if(this.props.info && _.size(this.props.info) > 0) {
      var machines = _.filter(this.props.info, function(machine) {
        return machine.Visible;
      });
      MachineNode = _.map(machines, function(machine) {
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
          <MachineChooser
            key={machine.Id}
            info={machine}
            user={this.props.user}
            busy={isMachineBusy}
            sameUser={isSameUser}
            activation={activationProps}
          />
        );
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
        {MachineNode}
      </div>
    );
  }
});

export default MachineList;

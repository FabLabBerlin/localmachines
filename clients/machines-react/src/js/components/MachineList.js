import React from 'react';
import MachineChooser from './MachineChooser';

/*
 * MachineList component:
 * Create a list of MachineChooser
 * Prepare the props for the MachineChooser to display the right MachineDiv
 */
var MachineList = React.createClass({

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
    if(this.props.info.length != 0) {
      var MachineNode = this.props.info.map(function(machine) {
        let activationProps = false;
        let isMachineBusy = false;
        let isSameUser = false;
        for( let i in activation ) {
          if( machine.Id == activation[i].MachineId ) {
            isMachineBusy = true;
            isSameUser = this.props.user.Id == activation[i].UserId;
            activationProps = activation[i];
            break;
          }
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
      var MachineNode = <p>You do not have access to any machines</p>;
    }
    return (
      <div className="machines">
        {MachineNode}
      </div>
    );
  }
});

module.exports = MachineList;

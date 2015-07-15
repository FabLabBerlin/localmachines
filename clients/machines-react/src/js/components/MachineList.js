import React from 'react';
import MachineChooser from './MachineChooser';

/*
 *  MachineList component:
 *  make a table where all the machine the user can access
 */
var MachineList = React.createClass({

  /*
   * Create the row of the table for each machine its get by props
   */
  render() {
    var activation = this.props.activation;
    var activationProps = false;
    if(this.props.info.length != 0) {
      var MachineNode = this.props.info.map(function(machine) {
        var isMachineBusy = false;
        var isSameUser = false;
        for( var i in activation ) {
          if( machine.Id == activation[i].MachineId ) {
            isMachineBusy = true;
            isSameUser = this.props.uid == activation[i].UserId;
            activationProps = activation[i];
            break;
          }
        }
        return (
          <MachineChooser
            key={machine.Id}
            info={machine}
            uid={this.props.uid}
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
      <div className="container-fluid" >
        {MachineNode}
      </div>
    );
  }
});

module.exports = MachineList;

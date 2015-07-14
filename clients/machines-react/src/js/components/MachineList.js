import React from 'react';
import Machine from './Machine';

/*
 *  MachineList component:
 *  make a table where all the machine the user can access
 */
var MachineList = React.createClass({

  /*
   * Create the row of the table for each machine its get by props
   */
  render() {
    if(this.props.info.length != 0) {
      var MachineNode = this.props.info.map(function(machine) {
        return (
          <Machine 
            key={machine.Id}
            info={machine}
            activation={this.props.activation}
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

import React from 'react';

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
          <tr key={machine.Id}>
            <td>{machine.Name}</td>
            <td>{machine.Shortname}</td>
            <td>{machine.Description}</td>
          </tr>
        );
      });
    } else {
      var MachineNode = <p>You do not have access to any machines</p>;
    }
    return (
      <table className="table table-striped table-hover" >
        <thead>
          <tr>
            <th>Machine Name</th>
            <th>Machine Shortname</th>
            <th>Description</th>
          </tr>
        </thead>
        <tbody>
          {MachineNode}
        </tbody>
      </table>
    );
  }
});

module.exports = MachineList;

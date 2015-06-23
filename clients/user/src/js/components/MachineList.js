import React from 'react';

var MachineList = React.createClass({
    render() {
        var MachineNode = this.props.info.map(function(machine) {
            return (
                <div key={machine.Id} >
                    <ul>
                        <li>Machine name: {machine.Name}</li>
                        <li>Description: {machine.Description}</li>
                    </ul>
                </div>
            );
        });
        return (
            <div className="machineList" >
                <p>Machine List</p>
                {MachineNode}
            </div>
        );
    }
});

module.exports = MachineList;

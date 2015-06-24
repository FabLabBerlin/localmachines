import React from 'react';

var MachineList = React.createClass({
    render() {
        if(this.props.info.length != 0) {
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
        } else {
            var MachineNode = <p>You do not have access to any machines</p>;
        }
        return (
            <div className="machineList" >
                <p>Machine List</p>
                {MachineNode}
            </div>
        );
    }
});

module.exports = MachineList;

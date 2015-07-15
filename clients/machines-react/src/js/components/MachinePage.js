import React from 'react'
import MachineList from './MachineList';

var MachinePage = React.createClass({
  getUserInfo() {
    return {
      uid: 3,
      FirstName: 'test',
      LastName: 'test'
    }
  },

  getAcInfo() {
    return [
      {
        Id: 1,
        UserId: 3,
        MachineId: 1,
        TimeStart: "2015-07-14T18:17:28+02:00"
      },
      {
        Id: 2,
        UserId: 2,
        MachineId: 3,
        TimeStart: "2015-07-14T18:17:28+02:00"
      }
    ]
  },

  getMachineInfo() {
    return [
      {
        Id: 1,
        Name: "Laydrop 3D Printer",
        Image: "",
        ConnectedMachines: "[3]",
        SwitchRefCount: 0
      },
      {
        Id: 3,
        Name: "Shitty Machine I dunno the name",
        Image: "",
        ConnectedMachines: "[2]",
        SwitchRefCount: 0
      },
      {
        Id: 2,
        Name: "MakerBot 3D Printer",
        Image: "machine-2.svg",
        ConnectedMachines: "",
        SwitchRefCount: 0
      }
    ]
  },

  getInitialState() {
    return {
      userInfo: this.getUserInfo(),
      machineInfo: this.getMachineInfo(),
      activationInfo: this.getAcInfo()
    };
  },

  getUserId() {
    return this.state.userInfo.uid;
  },

  render() {
    return (
      <div className="container-fluid" >
        <div>
          coucou {this.state.userInfo.FirstName} {this.state.userInfo.LastName}
        </div>
        <div>
          <MachineList 
            uid={this.getUserId()}
            info={this.state.machineInfo} 
            activation={this.state.activationInfo}
          />
        </div>
      </div>
    );
  }
});

module.exports = MachinePage;

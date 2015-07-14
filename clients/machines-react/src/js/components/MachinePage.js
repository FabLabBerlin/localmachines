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
    ]
  },

  getMachineInfo() {
    return [
      {
        Id: 1,
        Name: "Laydrop 3D Printer",
        Shortname: "MB3DP",
        Description: "NYC 3D printer 4 real and 4 life.",
        Image: "",
        Available: true,
        UnavailMsg: "",
        UnavailTill: "0001-01-01T00:00:00Z",
        Price: 16,
        PriceUnit: "hour",
        Comments: "",
        Visible: true,
        ConnectedMachines: "[3]",
        SwitchRefCount: 0
      },
      {
        Id: 2,
        Name: "MakerBot 3D Printer",
        Shortname: "MB3DP",
        Description: "NYC 3D printer 4 real and 4 life.",
        Image: "machine-2.svg",
        Available: true,
        UnavailMsg: "",
        UnavailTill: "0001-01-01T00:00:00Z",
        Price: 16,
        PriceUnit: "hour",
        Comments: "",
        Visible: true,
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

  render() {
    return (
      <div className="container-fluid" >
        <div>
          coucou {this.state.userInfo.FirstName} {this.state.userInfo.LastName}
        </div>
        <div>
          <MachineList 
            info={this.state.machineInfo} 
            activation={this.state.activationInfo}
          />
        </div>
      </div>
    );
  }
});

module.exports = MachinePage;

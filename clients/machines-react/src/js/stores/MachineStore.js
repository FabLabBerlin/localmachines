var MachineStore = {

  state: {
    userInfo: {
      uid: 3,
      FirstName: 'test',
      LastName: 'test'
    },
    activationInfo: [
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
    ],
    machineInfo: [
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

  endActivation(aid) {
    /*
    console.log('end activation');
    console.log(aid);
    */
  },

  startActivation(mid){
    /*
    console.log('start activation');
    console.log(mid);
    */
    var acNode = {
      Id: 5,
      UserId: 3,
      MachineId: mid,
      TimeStart: "blablablablablablabla"
    };
    this.state.activationInfo.push(acNode);
    this.onChangeActivation();
  },

  getUserInfo() {
    return this.state.userInfo;
  },

  getActivationInfo() {
    return this.state.activationInfo;
  },

  getMachineInfo() {
    return this.state.machineInfo;
  },

  onChangeActivation() {}

}

module.exports = MachineStore;

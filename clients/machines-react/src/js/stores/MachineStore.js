
/*
 * Import toastr and set position
 */
import toastr from 'toastr';
toastr.options.positionClass = 'toast-bottom-left';

var MachineStore = {

  state: {
    isLogged: false,
    firstTry: true,
    userInfo: {},
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

  /*
   * Login
   * submit the login form and try to connect to the back-end
   */
  submitLoginFormToServer(loginInfo){
    $.ajax({
      url: '/api/users/login',
      dataType: 'json',
      type: 'POST',
      data: loginInfo,
      success: function(data) {
        this.state.userInfo.UserId = data.UserId;
        console.log(this.state.userInfo);
        this.state.firstTry = true;
        this.state.isLogged = true;
        this.onChangeLogin();
      }.bind(this),
      error: function(xhr, status, err) {
        if(this.state.firstTry === true) {
          this.state.firstTry = false;
        } else {
          toastr.error('Wrong password or username');
        }
        console.error('/users/login', status, err.toString());
      }.bind(this),
    });
  },

  logout() {
    $.ajax({
      url: '/api/users/logout',
      type: 'GET',
      dataType: 'json',
      cache: false,
      success: function() {
        toastr.success('logout');
        this.state.isLogged = false;
        this.onChangeLogout();
      }.bind(this),
      error: function(xhr, status, err) {
        console.error('/users/logout', status, err.toString());
      }.bind(this)
    });
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

  getIsLogged() {
    return this.state.isLogged;
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

  onChangeLogin() {},

  onChangeLogout() {},

  onChangeActivation() {}

}

module.exports = MachineStore;

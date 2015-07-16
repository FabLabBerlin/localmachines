
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
   * POST call to the API
   * Make POST call cutomisable
   */
  postAPICall(url, dataToSend, successFunction, errorFunction) {
    $.ajax({
      url: url,
      dataType: 'json',
      type:'POST',
      data: dataToSend,
      success: function(data) {
        successFunction(data);
      }.bind(this),
      error: function(xhr, status, err) {
        errorFunction(xhr, status, err);
      }.bind(this),
    });
  },

  /*
   * GET call to the API
   * Make GET call cutomisable
   */
  getAPICall(url, successFunction, errorFunction) {
    $.ajax({
      url: url,
      dataType: 'json',
      type: 'GET',
      cache: false,
      success: function(data) {
        successFunction(data);
      }.bind(this),
      error: function(xhr, status, err) {
        errorFunction(xhr, status, err);
      }.bind(this),
    });
  },

  /*
   * Use POST call to login to the server
   * callback are defined below
   */
  submitLoginFormToServer(loginInfo){
    this.postAPICall('/api/users/login', loginInfo, this.LoginSuccessCallback, this.LoginErrorCallback);
  },

  /*
   * Activated when login succeed
   * MachineStore instead of this otherwise it doesn't work
   */
  LoginSuccessCallback(data) {
    MachineStore.state.userInfo.UserId = data.UserId;
    MachineStore.state.firstTry = true;
    MachineStore.state.isLogged = true;
    MachineStore.onChangeLogin();
  },

  /*
   * Activated when login failed
   * MachineStore instead of this otherwise it doesn't work
   */
  LoginErrorCallback(xhr, status, err) {
    if(MachineStore.state.firstTry === true) {
      MachineStore.state.firstTry = false;
    } else {
      toastr.error('Wrong password');
    }
    console.error('/users/login', status, err.toString());
  },

  /*
   * Use GET call to logout from the server
   * callback are defined below
   */
  logout() {
    getAPICall('/api/users/logou', logoutSuccessCallback, logoutErrorCallback);
  },

  /*
   * Callback for logout if succeed
   */
  logoutSuccessCallback(data) {
    toastr.success('logout');
    MachineStore.state.isLogged = false;
    MachineState.onChangeLogout();
  },

  /*
   * Callback for logout if failed
   */
  logoutErrorCallBack(xhr, status, err) {
    console.log('/users/logout', status, err.toString());
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

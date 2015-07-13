import $ from 'jquery';

/*
 * Import toastr and set the position
 */
import toastr from 'toastr';
toastr.options.positionClass = 'toast-bottom-left';

/*
 * @UserStore:
 * All the information about the user are stored here
 * All the interaction between the front-end and the back-end are done here
 * @state
 * @ajax calls:
 *  - PUT
 *  - POST
 *  - GET
 * @formatfunction
 * @getter
 * TODO: It would be possible to have multiple membership
 *       or to keep trace of your spending in membership
 *       Need to Change the reponse of /users/uid/membership by an array 
 *
 */
var UserStore = {
  /*
   * The state of the store
   */
  _state: {
    userId: 0,
    isLogged: false,
    firstTry: true,
    rawInfoUser: {},
    rawInfoMachine: [],
    rawInfoBill: {},
    rawInfoMembership: []
  },

  /*
   * Update the user profile
   * try to put the information in the UserForm into the database
   */
  submitUpdatedStateToServer(userState) {
    var updatedState = this.formatUserStateToSendToServer(userState);
    $.ajax({
      headers : {'Content-Type' : 'application/json'},
      url: '/api/users/' + this._state.userId,
      type: 'PUT',
      data: JSON.stringify( updatedState ),
      success: function() {
        toastr.success('Status updated');
      }.bind(this),
      error: function(xhr, status, err) {
        toastr.error('Error updating');
        console.error('/users/{uid}', status, err.toString());
      }.bind(this),
    });
  },

  /*
   * Login
   * submit the login form and try to connect to the back-end
   */
  submitLoginFormToServer(loginInfo) {
    $.ajax({
      url: '/api/users/login',
      dataType: 'json',
      type: 'POST',
      data: loginInfo,
      success: function(data) {
        this._state.userId = data.UserId;
        this.getUserStateFromServer(this._state.userId);
        this._state.firstTry = true;
      }.bind(this),
      error: function(xhr, status, err) {
        if(this._state.firstTry === true) {
          this._state.firstTry = false;
        } else {
          toastr.error('Wrong password or username');
        }
        console.error('/users/login', status, err.toString());
      }.bind(this),
    });
  },

  /*
   * Update the user's password
   */
  updatePassword(newPassword) {
    $.ajax({
      url: '/api/users/' + this._state.userId + '/password',
      dataType: 'json',
      type: 'POST',
      data: {
        password: newPassword
      },
      success: function() {
        toastr.success('Password successfully updated');
      }.bind(this),
      error: function(xhr, status, err) {
        toastr.error('Error while trying to update password');
        console.error('/users/{uid}/password', status, err.toString());
      }.bind(this),
    });
  },

  /*
   * Fetch User Data and store them
   * fetch the user info and call getMachineFromServer
   */
  getUserStateFromServer(uid){
    $.ajax({
      url: '/api/users/' + uid,
      dataType: 'json',
      type: 'GET',
      success: function(data) {
        this._state.rawInfoUser = data;
        this.getMachineFromServer(uid);
      }.bind(this),
      error: function(xhr, status, err) {
        toastr.error('Error getting the user\'s information');
        console.error('/users/{uid}', status, err.toString());
      }.bind(this),
    });
  },

  /*
   * Fetch Machines the user can use and store them
   * call getInfoBillFromServer if successful
   */
  getMachineFromServer(uid){
    $.ajax({
      url: '/api/users/' + uid + '/machinepermissions',
      dataType: 'json',
      type: 'GET',
      success: function(data) {
        this._state.rawInfoMachine = data;
        this.getInfoBillFromServer(uid);
      }.bind(this),
      error: function(xhr, status, err) {
        toastr.error('Error getting the user\'s machines');
        console.error('/users/{uid}/machinepermissions', status, err.toString());
      }.bind(this),
    });
  },

  /*
   * Fetch bill information and store them
   * call getMembershipFromServer if sucessful
   */
  getInfoBillFromServer(uid) {
    $.ajax({
      url: '/api/users/' + uid + '/bill',
      dataType: 'json',
      type: 'GET',
      success: function(data) {
        this._state.rawInfoBill = data;
        this.getMembershipFromServer(uid);
      }.bind(this),
      error: function(xhr, status, err) {
        toastr.error('Error getting the user\'s bill information');
        console.error('/users/{uid}/bill', status, err.toString());
      }.bind(this),
    });
  },

  /*
   * Fetch the membership the user subscribe and store it
   * call onChange if successful to alert UserPage
   */
  getMembershipFromServer(uid) {
    $.ajax({
      url: '/api/users/'+ uid +'/memberships',
      dataType: 'json',
      type: 'GET',
      success: function(data) {
        this._state.rawInfoMembership = data;
        this._state.isLogged = true;
        this.onChange();
      }.bind(this),
      error: function(xhr, status, err) {
        toastr.error('Error getting the membership');
        console.error('/users/{uid}/memberships', status, err.toString());
      }.bind(this),
    });
  },

  /*
   * Ask the server to logout
   * If successful, clean the state
   */
  logoutFromServer() {
    $.ajax({
      url: '/api/users/logout',
      type: 'GET',
      dataType: 'json',
      cache: false,
      success: function() {
        this.cleanState();
      }.bind(this),
      error: function(xhr, status, err) {
        console.error('/users/logout', status, err.toString())
      }.bind(this)
    });
  },

  /*
   * Update the store stat before sending the state to the server
   * @userState: the State you got from the UserForm
   */
  formatUserStateToSendToServer(userState) {
    for(var data in userState) {
      this._state.rawInfoUser[data] = userState[data];
    }
    var specialUpdateFormat = { User: this._state.rawInfoUser };
    return specialUpdateFormat;
  },

  /*
   * Create an Object with only the information we want to update and display
   */
  formatUserStateToSendToUserPage() {
    var infoWhichMatter = [
      'FirstName',
      'LastName', 
      'Username',
      'Email',
      'InvoiceAddr',
      'ShipAddr'
    ];
    var lightState = {};
    for(var index in infoWhichMatter) {
      var name = infoWhichMatter[index];
      lightState[name] = this._state.rawInfoUser[name];
    }
    return lightState;
  },

  /*
   * Clear the store
   */
  cleanState() {
    this._state.isLogged = false;
    this._state.userId = 0;
    this._state.rawInfoUser = {};
    this._state.rawInfoMachine = [];
    this._state.rawInfoBill = {};
    this._state.rawInfoMembership = {};
    this.onChangeLogout();
    toastr.success('Bye');
  },

  /*
   * Get the UID
   */
  getUID () {
    return this._state.userId;
  },

  /*
   * To know if the user is Logged or not
   */
  getIsLogged: function() {
    return this._state.isLogged;
  },

  /*
   * Return only the information which matter
   */
  getInfoUser: function() {
    var lightState = this.formatUserStateToSendToUserPage();
    return lightState;
  },

  /*
   * Return the machine the user can use
   */
  getInfoMachine: function() {
    return this._state.rawInfoMachine;
  },

  /*
   * return the detailled Bill of the user
   */
  getInfoBill: function() {
    return this._state.rawInfoBill;
  },

  /*
   * Return the membership the user have
   */
  getMembership() {
    return this._state.rawInfoMembership;
  },

  /*
   * Event when the user logout
   */
  onChangeLogout() {},

  /*
   * Event when all the data from the server are fetch
   */
  onChange() {}

};

module.exports = UserStore;

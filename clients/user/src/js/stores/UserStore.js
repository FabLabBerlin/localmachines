import $ from 'jquery';

/*
 * @UserStore:
 * All the information about the user are stored here
 * All the interaction between the front-end and the back-end are done here
 */
var UserStore = {
  /*
   * The state of the store
   */
  _state: {
    userID: 0,
    isLogged: false,
    rawInfoUser: {},
    rawInfoMachine: [],
    rawInfoMembership:{}
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
   * Update the user profile
   * try to put the information in the UserForm into the database
   */
  submitUpdatedStateToServer(userState) {
    var updatedState = this.formatUserStateToSendToServer(userState);
    $.ajax({
      headers : {'Content-Type' : 'application/json'},
      url: '/api/users/' + this._state.userID,
      type: 'PUT',
      data: JSON.stringify( updatedState ),
      success: function() {
        window.alert('change done');
      }.bind(this),
      error: function(xhr, status, err) {
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
        this._state.userID = data.UserId;
        this.getUserStateFromServer(this._state.userID);
      }.bind(this),
      error: function(xhr, status, err) {
        console.error('/users/login', status, err.toString());
        //invoke toaster stuff
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
        console.error('/users/{uid}', status, err.toString());
      }.bind(this),
    });
  },

  /*
   * Fetch Machine the user can use and store them
   * call getMembershipFromServer if successful
   */
  getMachineFromServer(uid){
    $.ajax({
      url: '/api/users/' + uid + '/machinepermissions',
      dataType: 'json',
      type: 'GET',
      success: function(data) {
        this._state.rawInfoMachine = data;
        this.getMembershipFromServer(uid);
      }.bind(this),
      error: function(xhr, status, err) {
        console.error('/users/{uid}/machinepermissions', status, err.toString());
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
        console.error('/users/{uid}/memberships', status, err.toString());
      }.bind(this),
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
    this._state.userID = 0;
    this._state.rawInfoUser = {};
    this._state.rawInfoMachine = [];
    this._state.rawInfoMembership = {};
    this.onChangeLogout();
  },

  /*
   * Get the UID
   */
  getUID () {
    return this._state.userID;
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

var $ = require('jquery');
var actionTypes = require('../actionTypes');
var ApiActions = require('./ApiActions');
var reactor = require('../reactor');
var toastr = require('../toastr');
var UserStore = require('../stores/UserStore');

/*
 * All the actions called by the UserPage
 */
var UserActions = {

  /*
   * Try to update the user information
   * @userState: data from userForm
   * call the UserStore to interact with the back-end
   */
  submitState(uid, userInfo){
    $.ajax({
      headers: {'Content-Type': 'application/json'},
      url: '/api/users/' + uid,
      type: 'PUT',
      data: JSON.stringify({
        User: userInfo
      }),
      success: function() {
        toastr.success('Status updated');
      }.bind(this),
      error: function(xhr, status, err) {
        toastr.error('Error updating');
        console.error('/users/{uid}', status, err.toString());
      }.bind(this)
    });
  },

  getUserInfoFromServer(uid) {
    ApiActions.getCall('/api/users/' + uid, _userInfoSuccess);
  },

  setUserInfoProperty({ key, value }) {
    reactor.dispatch(actionTypes.SET_USER_INFO_PROPERTY, { key, value });
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
        reactor.dispatch(actionTypes.SET_BILL_INFO, { data });
      }.bind(this),
      error: function(xhr, status, err) {
        toastr.error('Error getting the user\'s bill information');
        console.error('/users/{uid}/bill', status, err.toString());
      }.bind(this)
    });
  },

  /*
   * Fetch the membership the user subscribe and store it
   * call onChange if successful to alert UserPage
   */
  getMembershipFromServer(uid) {
    $.ajax({
      url: '/api/users/' + uid + '/memberships',
      dataType: 'json',
      type: 'GET',
      success: function(userMembershipList) {
        var data = userMembershipList.Data ? userMembershipList.Data : [];
        reactor.dispatch(actionTypes.SET_MEMBERSHIP_INFO, { data });
      }.bind(this),
      error: function(xhr, status, err) {
        toastr.error('Error getting the membership');
        console.error('/users/{uid}/memberships', status, err.toString());
      }.bind(this)
    });

  },

  /*
   * Ask the store to update the password
   * @password: new password the user want to have
   */
  updatePassword(uid, newPassword) {
    $.ajax({
      url: '/api/users/' + uid + '/password',
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
      }.bind(this)
    });
  }


};

/*
 * Success Callback
 * Activated when getNameLogin succeed
 * MachineStore instead of this otherwe it doesn't work
 */
function _userInfoSuccess(data) {
  var uid = data.Id;
  var usefulInformation = [
    'Id',
    'Username',
    'FirstName',
    'LastName',
    'Email',
    'Phone',
    'InvoiceAddr',
    'ZipCode',
    'City',
    'CountryCode',
    'UserRole',
    'Created',
    'Company'
  ];
  var userInfo = {};
  for(var index in usefulInformation) {
    userInfo[usefulInformation[index]] = data[usefulInformation[index]];
  }
  reactor.dispatch(actionTypes.SET_USER_INFO, { userInfo });
}

module.exports = UserActions;

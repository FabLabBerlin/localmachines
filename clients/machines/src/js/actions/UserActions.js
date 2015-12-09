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
  submitState(uid, user){
    $.ajax({
      headers: {'Content-Type': 'application/json'},
      url: '/api/users/' + uid,
      type: 'PUT',
      data: JSON.stringify({
        User: user
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

  fetchUser(uid) {
    ApiActions.getCall('/api/users/' + uid, function(user) {
      reactor.dispatch(actionTypes.SET_USER, { user });      
    });
  },

  setUserProperty({ key, value }) {
    reactor.dispatch(actionTypes.SET_USER_PROPERTY, { key, value });
  },

  /*
   * Fetch bill information and store them
   * call getMembershipFromServer if sucessful
   */
  fetchBill(uid) {
    $.ajax({
      url: '/api/users/' + uid + '/bill',
      dataType: 'json',
      type: 'GET',
      success: function(data) {
        reactor.dispatch(actionTypes.SET_BILL, { data });
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
  fetchMemberships(uid) {
    $.ajax({
      url: '/api/users/' + uid + '/memberships',
      dataType: 'json',
      type: 'GET',
      success(memberships) {
        var data = memberships.Data ? memberships.Data : [];
        reactor.dispatch(actionTypes.SET_MEMBERSHIPS, { data });
      },
      error(xhr, status, err) {
        toastr.error('Error getting the membership');
        console.error('/users/{uid}/memberships', status, err.toString());
      }
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

module.exports = UserActions;

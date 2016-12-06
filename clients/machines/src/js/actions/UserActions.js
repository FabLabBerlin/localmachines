var $ = require('jquery');
var actionTypes = require('../actionTypes');
var reactor = require('../reactor');
var toastr = require('../toastr');
var UserStore = require('../stores/UserStore');


var UserActions = {

  /*
   * Try to update the user information
   * @userState: data from userForm
   * call the UserStore to interact with the back-end
   */
  updateUser(uid, user){
    $.ajax({
      headers: {'Content-Type': 'application/json'},
      url: '/api/users/' + uid,
      type: 'PUT',
      data: JSON.stringify({
        User: user
      }),
      success() {
        toastr.success('Status updated');
      },
      error(xhr, status, err) {
        toastr.error('Error updating');
        console.error('/users/{uid}', status, err.toString());
      }
    });
  },

  fetchUser(uid) {
    $.ajax({
      url: '/api/users/' + uid
    })
    .done(user => {
      reactor.dispatch(actionTypes.SET_USER, { user });      
    });
  },

  setUserProperty({ key, value }) {
    reactor.dispatch(actionTypes.SET_USER_PROPERTY, { key, value });
  },

  fetchBill(locationId, uid) {
    $.ajax({
      url: '/api/users/' + uid + '/bill?location=' + locationId,
      dataType: 'json',
      type: 'GET',
      success(data) {
        reactor.dispatch(actionTypes.SET_BILL, { data });
      },
      error(xhr, status, err) {
        toastr.error('Error getting the user\'s bill information');
        console.error('/users/{uid}/bill', status, err.toString());
      }
    });
  },

  fetchMemberships(locationId, uid) {
    $.ajax({
      url: '/api/users/' + uid + '/memberships?location=' + locationId,
      dataType: 'json',
      type: 'GET',
      success(memberships) {
        var data = memberships.Data;
        reactor.dispatch(actionTypes.SET_MEMBERSHIPS, { data });
      },
      error(xhr, status, err) {
        toastr.error('Error getting the membership');
        console.error('/users/{uid}/memberships', status, err.toString());
      }
    });

  },

  updatePassword(uid, newPassword) {
    $.ajax({
      url: '/api/users/' + uid + '/password',
      dataType: 'json',
      type: 'POST',
      data: {
        password: newPassword
      },
      success() {
        toastr.success('Password successfully updated');
      },
      error(xhr, status, err) {
        toastr.error('Error while trying to update password');
        console.error('/users/{uid}/password', status, err.toString());
      }
    });
  }


};

export default UserActions;

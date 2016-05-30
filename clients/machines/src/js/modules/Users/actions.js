var $ = require('jquery');
var actionTypes = require('./actionTypes');
var reactor = require('../../reactor');
var toastr = require('../../toastr');


function fetchUsers({locationId}) {
  $.ajax({
    url: '/api/users?location=' + locationId,
    dataType: 'json'
  })
  .success((users) => {
    reactor.dispatch(actionTypes.SET_USERS, users);
  })
  .error(() => {
    toastr.error('Error getting users');
  });
}

export default {
  fetchUsers
};

var _ = require('lodash');
var $ = require('jquery');
var actionTypes = require('./actionTypes');
var reactor = require('../../reactor');
var toastr = require('../../toastr');


function fetchUserMemberships({locationId, userId}) {
  $.ajax({
    url: '/api/users/' + userId + '/bill?location=' + locationId,
    dataType: 'json'
  })
  .success((invoices) => {
    var userMemberships = [];

    _.each(invoices, function(invoice) {
      _.each(invoice.UserMemberships.Data, function(umb) {
        umb.Invoice = invoice;
        userMemberships.push(umb);
      });
    });

    reactor.dispatch(actionTypes.SET_USER_MEMBERSHIPS, {
      userId: userId,
      userMemberships: userMemberships
    });
  })
  .error(() => {
    toastr.error('Error getting user memberships');
  });
}

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
  fetchUsers,
  fetchUserMemberships
};

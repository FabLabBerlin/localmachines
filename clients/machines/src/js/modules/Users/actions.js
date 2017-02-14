import _ from 'lodash';
var $ = require('jquery');
import actionTypes from './actionTypes';
import moment from 'moment';
import reactor from '../../reactor';
import toastr from '../../toastr';


function addUserMembership({locationId, userId, membershipId}) {
  $.ajax({
    method: 'POST',
    url: '/api/users/' + userId + '/memberships',
    data: {
      membershipId: membershipId,
      startDate: moment().format('YYYY-MM-DD')
    }
  })
  .success(() => {
    toastr.success('Membership created');
    fetchUserMemberships({locationId, userId});
  })
  .error(() => {
    toastr.error('Error while trying to create new User Membership');
  });
}

function fetchMemberships({locationId}) {
  $.ajax({
    url: '/api/memberships?location=' + locationId,
    dataType: 'json'
  })
  .success((memberships) => {
    reactor.dispatch(actionTypes.SET_MEMBERSHIPS, {
      memberships: memberships
    });
  })
  .error(() => {
    toastr.error('Error getting memberships');
  });
}

function fetchUserMemberships({locationId, userId}) {
  $.ajax({
    url: '/api/users/' + userId + '/bill?location=' + locationId,
    dataType: 'json'
  })
  .success((invoices) => {
    var userMemberships = [];

    _.each(invoices, (invoice) => {
      _.each(invoice.UserMemberships.Data, (umb) => {
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
  addUserMembership,
  fetchMemberships,
  fetchUsers,
  fetchUserMemberships
};

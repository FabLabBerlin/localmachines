var $ = require('jquery');
var actionTypes = require('./actionTypes');
var reactor = require('../../reactor');
var toastr = require('../../toastr');


function fetchUserMemberships(locId, {userId}) {
  $.ajax({
    method: 'GET',
    url: '/api/users/' + userId + '/memberships',
    data: {
      location: locId
    }
  })
  .success(function(userMemberships) {
    reactor.dispatch(actionTypes.SET_USER_MEMBERSHIPS, {
      userId: userId,
      userMemberships: userMemberships
    });
  })
  .error(function() {
    toastr.error('Error fetch user memberships');
  });
}

function fetchMonthlySums(locId, {month, year}) {
  $.ajax({
    method: 'GET',
    url: '/api/invoices/months/' + year + '/' + month,
    data: {
      location: locId
    }
  })
  .success(function(summaries) {
    reactor.dispatch(actionTypes.FETCH_MONTHLY_SUMMARIES, {
      year: year,
      month: month,
      summaries: summaries
    });
  })
  .error(function() {
    toastr.error('Error fetch monthly summaries.  Please try again later.');
  });
}

function fetchUser(locId, {month, year, userId}) {
  $.ajax({
    method: 'GET',
    url: '/api/invoices/months/' + year + '/' + month + '/users/' + userId,
    data: {
      location: locId
    }
  })
  .success(function(invoice) {
    reactor.dispatch(actionTypes.SET_INVOICE, {
      year: year,
      month: month,
      userId: userId,
      invoice: invoice
    });
  })
  .error(function() {
    toastr.error('Error fetch monthly summaries.  Please try again later.');
  });
}

function selectUserId(userId) {
  reactor.dispatch(actionTypes.SELECT_USER_ID, userId);
}

function setSelectedMonth({month, year}) {
  reactor.dispatch(actionTypes.SET_SELECTED_MONTH, { month, year });
}

export default {
  fetchMonthlySums,
  fetchUser,
  fetchUserMemberships,
  selectUserId,
  setSelectedMonth
};

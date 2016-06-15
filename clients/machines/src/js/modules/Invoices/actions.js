var $ = require('jquery');
var actionTypes = require('./actionTypes');
var reactor = require('../../reactor');
var toastr = require('../../toastr');

function editPurchase(id) {
  reactor.dispatch(actionTypes.EDIT_PURCHASE, id);
}

function fetchFastbillStatuses(locId, {month, year, userId}) {
  $.ajax({
    url: '/api/billing/months/' + year + '/' + month + '/users/' + userId + '/statuses',
    data: {
      location: locId
    }
  })
  .success(function(invoiceStatuses) {
    console.log('invoiceStatuses=', invoiceStatuses);
    reactor.dispatch(actionTypes.SET_INVOICE_STATUSES, {
      month: month,
      year: year,
      userId: userId,
      invoiceStatuses: invoiceStatuses
    });
  })
  .error(function() {
    toastr.error('Error fetching invoice statuses.');
  });
}

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
  const selected = 

  $.ajax({
    method: 'GET',
    url: '/api/billing/months/' + year + '/' + month,
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

function fetchInvoice(locId, {invoiceId}) {
  $.ajax({
    method: 'GET',
    url: '/api/billing/invoices/' + invoiceId,
    data: {
      location: locId
    }
  })
  .success(function(invoice) {
    reactor.dispatch(actionTypes.SET_INVOICE, {
      invoice: invoice
    });
  })
  .error(function() {
    toastr.error('Error fetch monthly summaries.  Please try again later.');
  });
}

function makeDraft(locId, {month, year, userId}) {
  $.ajax({
    method: 'POST',
    url: '/api/billing/months/' + year + '/' + month + '/users/' + userId + '/draft',
    data: {
      location: locId
    }
  })
  .success(function(invoice) {
    console.log('invoice=', invoice);
    toastr.info('Draft created');
  })
  .error(function() {
    toastr.error('Error creating draft.');
  });
}

function save(locId, {month, year, userId}) {
  $.ajax({
    method: 'POST',
    url: '/api/billing/months/' + year + '/' + month + '/users/' + userId + '/update',
    data: {
      location: locId
    }
  })
  .success(function(invoice) {
    toastr.info('Changes saved');
  })
  .error(function() {
    toastr.error('Error saving changes.');
  });
}

function selectInvoiceId(invoiceId) {
  reactor.dispatch(actionTypes.SELECT_INVOICE_ID, invoiceId);
}

function setSelectedMonth({month, year}) {
  reactor.dispatch(actionTypes.SET_SELECTED_MONTH, { month, year });
}

export default {
  editPurchase,
  fetchFastbillStatuses,
  fetchInvoice,
  fetchMonthlySums,
  fetchUserMemberships,
  makeDraft,
  save,
  selectInvoiceId,
  setSelectedMonth
};

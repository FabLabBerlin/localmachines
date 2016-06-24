var _ = require('lodash');
var $ = require('jquery');
var actionTypes = require('./actionTypes');
var getters = require('./getters');
var reactor = require('../../reactor');
var toastr = require('../../toastr');

function complete() {
  /*eslint-disable no-alert */
  if (!window.confirm('Invoice cannot be changed and will be synchronized with Fastbill.')) {
    toastr.warning('Aborted complete invoice');
    return;
  }
  /*eslint-enable no-alert */

  const invoice = reactor.evaluateToJS(getters.getInvoice);

  $.ajax({
    method: 'POST',
    url: '/api/billing/invoices/' + invoice.Id + '/complete',
    data: {
      location: invoice.LocationId
    }
  })
  .success(function() {
    toastr.info('Invoice completed');
  })
  .error(function() {
    toastr.error('Error completing invoice.');
  });
}

function editPurchase(id) {
  reactor.dispatch(actionTypes.EDIT_PURCHASE, id);
}

function editPurchaseDuration(duration) {
  reactor.dispatch(actionTypes.EDIT_PURCHASE_DURATION, duration);
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

function makeDraft(locId) {
  const invoice = reactor.evaluateToJS(getters.getInvoice);

  if (invoice.FastbillId) {
    /*eslint-disable no-alert */
    if (!window.confirm('Invoice already pushed to Fastbill. Overwrite changes in Fastbill?')) {
      toastr.warning('Aborted make draft');
      return;
    }
    /*eslint-enable no-alert */
  }

  $.ajax({
    method: 'POST',
    url: '/api/billing/invoices/' + invoice.Id + '/draft',
    data: {
      location: locId
    }
  })
  .success(function() {
    toastr.info('Draft created');
  })
  .error(function() {
    toastr.error('Error creating draft.');
  });
}

function save(locId, {invoiceId}) {
  var inv = reactor.evaluateToJS(getters.getInvoice);
  var falseEdits = false;
  var mutated = _.filter(inv.Purchases, (p) => {
    if (p.editValid === false) {
      falseEdits = true;
    }

    return p.editedDuration;
  });

  if (falseEdits) {
    toastr.error('Trying to save invalid edit');
    return;
  }

  var promises = _.map(mutated, (p) => {
    return $.ajax({
      headers: {'Content-Type': 'application/json'},
      method: 'PUT',
      url: '/api/activations/' + p.Id,
      data: JSON.stringify(p),
      params: {
        location: locId
      }
    });
  });

  $.when(promises)
  .done(() => {
    toastr.info('Successfully saved updates');
    editPurchase(undefined);
    fetchInvoice(inv.LocationId, {
      invoiceId: inv.Id
    });
  })
  .fail(() => {
    toastr.error('Error while saving.');
  });
}

function selectInvoiceId(invoiceId) {
  reactor.dispatch(actionTypes.SELECT_INVOICE_ID, invoiceId);
}

function setSelectedMonth({month, year}) {
  reactor.dispatch(actionTypes.SET_SELECTED_MONTH, { month, year });
}

export default {
  complete,
  editPurchase,
  editPurchaseDuration,
  fetchFastbillStatuses,
  fetchInvoice,
  fetchMonthlySums,
  fetchUserMemberships,
  makeDraft,
  save,
  selectInvoiceId,
  setSelectedMonth
};

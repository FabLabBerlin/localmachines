var _ = require('lodash');
var $ = require('jquery');
var actionTypes = require('./actionTypes');
var getters = require('./getters');
var GlobalActions = require('../../actions/GlobalActions');
var reactor = require('../../reactor');
var toastr = require('../../toastr');

function cancel() {
  /*eslint-disable no-alert */
  if (!window.confirm('Really cancel invoice?')) {
    toastr.warning('Aborted canceling invoice');
    return;
  }
  /*eslint-enable no-alert */

  const invoice = reactor.evaluateToJS(getters.getInvoice);

  GlobalActions.showGlobalLoader();

  $.ajax({
    method: 'POST',
    url: '/api/billing/invoices/' + invoice.Id + '/cancel',
    data: {
      location: invoice.LocationId
    }
  })
  .success(() => {
    toastr.info('Invoice canceled');
  })
  .error(() => {
    toastr.error('Error canceling invoice.');
  })
  .always(() => {
    GlobalActions.hideGlobalLoader();
  });
}

function check(invoiceId) {
  reactor.dispatch(actionTypes.CHECK, invoiceId);
}

function checkAll() {
  reactor.dispatch(actionTypes.CHECK_ALL);
}

function checkedComplete() {
}

function checkedPushDrafts(locId) {
  var iterate = function(ids) {
    const id = ids.shift();

    $.ajax({
      method: 'POST',
      url: '/api/billing/invoices/' + id + '/draft',
      data: {
        location: locId
      }
    })
    .success(() => {
      toastr.info('Draft pushed');
      if (ids.length > 0) {
        iterate(ids);
      }
    })
    .error((jqXHR, textStatus) => {
      console.log(jqXHR.responseText);
      console.log('textStatus=', textStatus);
      toastr.error('Error pushing draft, aborting.');
      GlobalActions.hideGlobalLoader();
    })
    .always(() => {
      if (ids.length === 0) {
        GlobalActions.hideGlobalLoader();
      }
    });
  };

  const invs = _.filter(
    reactor.evaluateToJS(getters.getThisMonthInvoices),
    inv => inv.checked
  );

  if (invs.length > 0) {
    toastr.info('Will Push selected invoices in draft status');
    GlobalActions.showGlobalLoader();
    iterate(_.pluck(invs, 'Id'));
  } else {
    toastr.error('No invoices selected');
  }
}

function checkedSend() {
}

function checkSetStatus(status) {
  reactor.dispatch(actionTypes.CHECK_SET_STATUS, status);
}

function complete() {
  /*eslint-disable no-alert */
  if (!window.confirm('Invoice cannot be changed and will be synchronized with Fastbill.')) {
    toastr.warning('Aborted complete invoice');
    return;
  }
  /*eslint-enable no-alert */

  const invoice = reactor.evaluateToJS(getters.getInvoice);

  GlobalActions.showGlobalLoader();

  $.ajax({
    method: 'POST',
    url: '/api/billing/invoices/' + invoice.Id + '/complete',
    data: {
      location: invoice.LocationId
    }
  })
  .success(() => {
    toastr.info('Invoice completed');
  })
  .error(() => {
    toastr.error('Error completing invoice.');
  })
  .always(() => {
    GlobalActions.hideGlobalLoader();
    editPurchase(undefined);
    fetchInvoice(invoice.LocationId, {
      invoiceId: invoice.Id
    });
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
  GlobalActions.showGlobalLoader();

  $.ajax({
    method: 'GET',
    url: '/api/billing/invoices/' + invoiceId,
    data: {
      location: locId
    }
  })
  .success((invoice) => {
    reactor.dispatch(actionTypes.SET_INVOICE, {
      invoice: invoice
    });
  })
  .error(() => {
    toastr.error('Error fetch monthly summaries.  Please try again later.');
  })
  .always(() => {
    GlobalActions.hideGlobalLoader();
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

  GlobalActions.showGlobalLoader();

  $.ajax({
    method: 'POST',
    url: '/api/billing/invoices/' + invoice.Id + '/draft',
    data: {
      location: locId
    }
  })
  .success(() => {
    toastr.info('Draft created');
  })
  .error(() => {
    toastr.error('Error creating draft.');
  })
  .always(() => {
    GlobalActions.hideGlobalLoader();
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

  GlobalActions.showGlobalLoader();

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
  })
  .always(() => {
    GlobalActions.hideGlobalLoader();
  });
}

function selectInvoiceId(invoiceId) {
  reactor.dispatch(actionTypes.SELECT_INVOICE_ID, invoiceId);
}

function send() {
  /*eslint-disable no-alert */
  if (!window.confirm('Really send invoice?')) {
    toastr.warning('Aborted send invoice');
    return;
  }
  /*eslint-enable no-alert */

  const invoice = reactor.evaluateToJS(getters.getInvoice);

  GlobalActions.showGlobalLoader();

  $.ajax({
    method: 'POST',
    url: '/api/billing/invoices/' + invoice.Id + '/send',
    data: {
      location: invoice.LocationId
    }
  })
  .success(() => {
    toastr.info('Invoice sent');
  })
  .error(() => {
    toastr.error('Error sending invoice.');
  })
  .always(() => {
    GlobalActions.hideGlobalLoader();
    editPurchase(undefined);
    fetchInvoice(invoice.LocationId, {
      invoiceId: invoice.Id
    });
  });
}

function setSelectedMonth({month, year}) {
  reactor.dispatch(actionTypes.SET_SELECTED_MONTH, { month, year });
}

export default {
  cancel,
  check,
  checkAll,
  checkedPushDrafts,
  checkSetStatus,
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
  send,
  setSelectedMonth
};

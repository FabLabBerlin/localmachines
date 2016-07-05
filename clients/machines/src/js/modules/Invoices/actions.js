var _ = require('lodash');
var $ = require('jquery');
var actionTypes = require('./actionTypes');
var getters = require('./getters');
var GlobalActions = require('../../actions/GlobalActions');
var LocationGetters = require('../../modules/Location/getters');
var React = require('react');
var reactor = require('../../reactor');
var toastr = require('../../toastr');

// https://github.com/HubSpot/vex/issues/72
var vex = require('vex-js'),
VexDialog = require('vex-js/js/vex.dialog.js');

vex.defaultOptions.className = 'vex-theme-custom';


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
    refresh();
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

function checkedBatch(uriAction, verb) {
  const locId = reactor.evaluateToJS(LocationGetters.getLocationId);
  const ivs = _.filter(
    reactor.evaluateToJS(getters.getThisMonthInvoices),
    inv => inv.checked
  );
  console.log('ivs=', ivs);
  const n = ivs.length;

  var iterate = function(invs, errors) {
    const inv = invs.shift();
    const i = n - invs.length;

    if (i % 20 === 0) {
      toastr.info('Processed ' + i + '/' + n + ' invoices');
    }

    $.ajax({
      method: 'POST',
      url: '/api/billing/invoices/' + inv.Id + '/' + uriAction,
      data: {
        location: locId
      }
    })
    .success(() => {
      console.log(verb);
      if (invs.length > 0) {
        iterate(invs, errors);
      } else {
        finish(errors);
        refresh();
      }
    })
    .error((jqXHR, textStatus) => {
      const err = jqXHR.responseText;
      console.log(err);
      console.log('textStatus=', textStatus);
      console.log('Error ' + verb + '.');
      errors.push({
        message: err,
        user: inv.User
      });
      if (invs.length > 0) {
        iterate(invs, errors);
      } else {
        finish(errors);
      }
    });
  };

  var finish = function(errors) {
    GlobalActions.hideGlobalLoader();

    var msg = $('<div></div>');

    msg.append('<h2>Finished</h2>');

    if (errors.length > 0) {
      msg.append('<h4>Some errors occurred</h4>');
      _.each(errors, (err) => {
        const name = err.user.FirstName + ' ' + err.user.LastName;

        msg.append('<h6>Invoice of ' + name + ': ' + err.message + '</h6>');
      });
    } else {
      msg.append('<h4>Successful ' + verb + ' for all</h4>');
    }
    VexDialog.alert(msg.html());
  };

  if (ivs.length > 0) {
    toastr.info('Selected invoices: ' + verb);
    GlobalActions.showGlobalLoader();
    iterate(ivs, []);
  } else {
    toastr.error('No invoices selected');
  }  
}

function checkedComplete() {
  checkedBatch('complete', 'freeze');
}

function checkedPushDrafts() {
  checkedBatch('draft', 'draft-pushing');
}

function checkedSend() {
  checkedBatch('send', 'send');
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
    refresh();
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

function refresh() {
  const inv = reactor.evaluateToJS(getters.getInvoice);
  const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
  const monthlySums = reactor.evaluateToJS(getters.getMonthlySums);
  const month = monthlySums.selected.month;
  const year = monthlySums.selected.year;

  if (inv) {
    fetchInvoice(inv.LocationId, {
      invoiceId: inv.Id
    });
  }
  fetchMonthlySums(locationId, {
    month: month,
    year: year
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
    refresh();
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

function send(canceled) {
  const msg = canceled ? 'Really send cancelation invoice?' : 'Really send invoice?';

  /*eslint-disable no-alert */
  if (!window.confirm(msg)) {
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
      location: invoice.LocationId,
      canceled: canceled
    }
  })
  .success(() => {
    if (canceled) {
      toastr.info('Cancelation invoice sent');
    } else {
      toastr.info('Invoice sent');
    }
  })
  .error(() => {
    if (canceled) {
      toastr.error('Error sending cancelation invoice.');
    } else {
      toastr.error('Error sending invoice.');
    }
  })
  .always(() => {
    GlobalActions.hideGlobalLoader();
    editPurchase(undefined);
    refresh();
  });
}

function sendCanceled() {
  send(true);
}

function setSelectedMonth({month, year}) {
  reactor.dispatch(actionTypes.SET_SELECTED_MONTH, { month, year });
}

function sortBy(column, asc) {
  reactor.dispatch(actionTypes.SORT_BY, {column, asc});
}

export default {
  cancel,
  check,
  checkAll,
  checkedComplete,
  checkedPushDrafts,
  checkedSend,
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
  sendCanceled,
  setSelectedMonth,
  sortBy
};

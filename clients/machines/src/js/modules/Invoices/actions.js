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


function cancel(invoice) {
  /*eslint-disable no-alert */
  if (!window.confirm('Really cancel invoice?')) {
    toastr.warning('Aborted canceling invoice');
    return;
  }
  /*eslint-enable no-alert */

  GlobalActions.showGlobalLoader();

  $.ajax({
    method: 'POST',
    url: '/api/billing/invoices/' + invoice.get('Id') + '/cancel',
    data: {
      location: invoice.get('LocationId')
    }
  })
  .success(() => {
    toastr.info('Invoice canceled');
    refresh(invoice);
  })
  .error((xhr) => {
    toastr.error('Error canceling invoice:' + xhr.responseText);
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
  /*eslint-disable no-alert */
  if (!window.confirm('Really ' + verb + ' selected invoices?')) {
    toastr.warning('Aborted action');
    return;
  }
  /*eslint-enable no-alert */

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
        refresh(inv);
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

function complete(invoice) {
  /*eslint-disable no-alert */
  if (!window.confirm('Invoice cannot be changed and will be synchronized with Fastbill.')) {
    toastr.warning('Aborted complete invoice');
    return;
  }
  /*eslint-enable no-alert */

  GlobalActions.showGlobalLoader();

  $.ajax({
    method: 'POST',
    url: '/api/billing/invoices/' + invoice.get('Id') + '/complete',
    data: {
      location: invoice.get('LocationId')
    }
  })
  .success(() => {
    toastr.info('Invoice completed');
  })
  .error((xhr) => {
    toastr.error('Error completing invoice:' + xhr.responseText);
  })
  .always(() => {
    GlobalActions.hideGlobalLoader();
    editPurchase(undefined);
    refresh(invoice);
  });
}

function createPurchase(invoice) {
  console.log('invoice=', invoice);
  GlobalActions.showGlobalLoader();

  $.ajax({
    method: 'POST',
    url: '/api/activations',
    data: {
      location: invoice.get('LocationId'),
      invoice: invoice.get('Id'),
      user: invoice.get('UserId')
    }
  })
  .success((data) => {
    toastr.info('Purchase created');
    console.log('data=', data);
    refresh(invoice);
    editPurchase(data.Id);
  })
  .error(xhr => {
    toastr.error('Error creating purchase:' + xhr.responseText);
  })
  .always(() => {
    GlobalActions.hideGlobalLoader();
  });
}

function editPurchase(id) {
  reactor.dispatch(actionTypes.EDIT_PURCHASE, id);
}

function editPurchaseDuration(invoice, duration) {
  const invoiceId = invoice.get('Id');

  reactor.dispatch(actionTypes.EDIT_PURCHASE_DURATION, {duration, invoiceId});
}

function editPurchaseField({invoice, field, value}) {
  const invoiceId = invoice.get('Id');

  reactor.dispatch(actionTypes.EDIT_PURCHASE_FIELD, {invoiceId, field, value});
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
  .success((invoice) => {
    reactor.dispatch(actionTypes.SET_INVOICE, {
      invoice: invoice
    });
  })
  .error(() => {
    toastr.error('Error fetch monthly summaries.  Please try again later.');
  });
}

function makeDraft(locId, invoice) {
  if (invoice.get('FastbillId')) {
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
    url: '/api/billing/invoices/' + invoice.get('Id') + '/draft',
    data: {
      location: locId
    }
  })
  .success(() => {
    toastr.info('Draft created');
  })
  .error((xhr) => {
    toastr.error('Error creating draft:' + xhr.responseText);
  })
  .always(() => {
    GlobalActions.hideGlobalLoader();
  });
}

function refresh(inv) {
  const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
  const monthlySums = reactor.evaluateToJS(getters.getMonthlySums);
  const month = monthlySums.selected.month;
  const year = monthlySums.selected.year;

  if (inv) {
    fetchInvoice(inv.get('LocationId'), {
      invoiceId: inv.get('Id')
    });
  }
  fetchMonthlySums(locationId, {
    month: month,
    year: year
  });
}

function save(locId, {invoice}) {
  console.log('Invoice actions#save');
  var falseEdits = false;

  var mutated = _.filter(invoice.get('Purchases').toJS(), (p) => {
    if (p.editValid === false) {
      falseEdits = true;
    }

    return p.edited || p.editedDuration;
  });

  console.log('falseEdits=', falseEdits);

  if (falseEdits) {
    toastr.error('Trying to save invalid edit');
    return;
  }

  var promises = _.map(mutated, (p) => {
    var url;

    switch (p.Type) {
    case 'activation':
      url = '/api/activations/' + p.Id + '?location=' + locId;
      break;
    case 'reservation':
      url = '/api/reservations/' + p.Id + '?location=' + locId;
      break;
    default:
      url = '/api/purchases/' + p.Id + '?type=' + p.Type;
    }

    console.log('sending data=', p);

    if (p.Type === 'other') {
      if (!_.isNumber(p.PricePerUnit)) {
        p.PricePerUnit = parseFloat(p.PricePerUnit);
      }

      if (!_.isNumber(p.Quantity)) {
        p.Quantity = parseFloat(p.Quantity);
      }
    } else {
      if (!_.isNumber(p.MachineId)) {
        p.MachineId = parseInt(p.MachineId);
      }
    }

    return $.ajax({
      headers: {'Content-Type': 'application/json'},
      method: 'PUT',
      url: url,
      data: JSON.stringify(p),
      params: {
        location: locId
      }
    });
  });

  GlobalActions.showGlobalLoader();

  $.when.apply(this, promises)
  .done((...results) => {
    toastr.info('Successfully saved updates');
    editPurchase(undefined);
    refresh(invoice);
  })
  .fail((xhr) => {
    toastr.error('Error while saving:' + xhr.responseText);
  })
  .always(() => {
    GlobalActions.hideGlobalLoader();
  });
}

function send(canceled, invoice) {
  const msg = canceled ? 'Really send cancelation invoice?' : 'Really send invoice?';

  /*eslint-disable no-alert */
  if (!window.confirm(msg)) {
    toastr.warning('Aborted send invoice');
    return;
  }
  /*eslint-enable no-alert */

  GlobalActions.showGlobalLoader();

  $.ajax({
    method: 'POST',
    url: '/api/billing/invoices/' + invoice.get('Id') + '/send',
    data: {
      location: invoice.get('LocationId'),
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
  .error((xhr) => {
    if (canceled) {
      toastr.error('Error sending cancelation invoice:' + xhr.responseText);
    } else {
      toastr.error('Error sending invoice:' + xhr.responseText);
    }
  })
  .always(() => {
    GlobalActions.hideGlobalLoader();
    editPurchase(undefined);
    refresh(invoice);
  });
}

function sendCanceled(invoice) {
  send(true, invoice);
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
  createPurchase,
  editPurchase,
  editPurchaseDuration,
  editPurchaseField,
  fetchFastbillStatuses,
  fetchInvoice,
  fetchMonthlySums,
  fetchUserMemberships,
  makeDraft,
  save,
  send,
  sendCanceled,
  setSelectedMonth,
  sortBy
};

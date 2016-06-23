var _ = require('lodash');
var actionTypes = require('../actionTypes');
var moment = require('moment');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;

const initialState = toImmutable({
  MonthlySums: {
  invoices: {},
  selected: {
      month: moment().month() + 1,
      year: moment().year(),
      userId: undefined
    }
  }
});


var InvoicesStore = new Nuclear.Store({

  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.EDIT_PURCHASE, editPurchase);
    this.on(actionTypes.EDIT_PURCHASE_DURATION, editPurchaseDuration);
    this.on(actionTypes.FETCH_MONTHLY_SUMMARIES, fetchMonthlySums);
    this.on(actionTypes.SELECT_INVOICE_ID, selectInvoiceId);
    this.on(actionTypes.SET_SELECTED_MONTH, setSelectedMonth);
    this.on(actionTypes.SET_INVOICE, setInvoice);
    this.on(actionTypes.SET_INVOICE_STATUSES, setInvoiceStatuses);
    this.on(actionTypes.SET_USER_MEMBERSHIPS, setUserMemberships);
  }

});

function editPurchase(state, id) {
  return state.set('editPurchaseId', id);
}

function editPurchaseDuration(state, duration) {
  var invoiceId = state.getIn(['MonthlySums', 'selected', 'invoiceId']);
  var purchaseId = state.get('editPurchaseId');

  var keyPath = [
    'invoices',
    'detailedInvoices',
    invoiceId
  ];

  var iv = state.getIn(keyPath).toJS();
  iv.Purchases = _.map(iv.Purchases, (p) => {
    if (p.Id === purchaseId) {
      p.editValid = false;
      p.editedDuration = duration;

      if (duration.length === 10) {
        var timeEnd = toTimeEnd(p, duration);
        if (timeEnd) {
          p.editValid = true;
          p.TimeEnd = timeEnd.toDate();
        }
      }
    }
    return p;
  });

  return state.setIn(keyPath, toImmutable(iv));
}

function fetchMonthlySums(state, { month, year, summaries }) {
  return state.setIn(['MonthlySums', year, month], toImmutable(summaries));
}

function selectInvoiceId(state, invoiceId) {
  return state.setIn(['MonthlySums', 'selected', 'invoiceId'], invoiceId);
}

function setSelectedMonth(state, { month, year }) {
  return state.setIn(['MonthlySums', 'selected', 'month'], month)
              .setIn(['MonthlySums', 'selected', 'year'], year);
}

function setInvoice(state, { invoice }) {
  return state.setIn(['invoices', 'detailedInvoices', invoice.Id], toImmutable(invoice));
}

function setInvoiceStatuses(state, { month, year, userId, invoiceStatuses }) {
  return state.setIn(['invoiceStatuses', year, month, userId], invoiceStatuses);
}

function setUserMemberships(state, { userId, userMemberships }) {
  return state.setIn(['userMemberships', userId], userMemberships);
}

// Private:

function toTimeEnd(p, duration) {
  var str = duration.slice(0, 8);
  var hms = str.split(':');

  if (hms.length !== 3) {
    return undefined;
  }

  var hh = parseInt(hms[0], 10);
  var mm = parseInt(hms[1], 10);
  var ss = parseInt(hms[2], 10);

  if (!_.isNumber(hh) || !_.isNumber(mm) || !_.isNumber(ss)) {
    return undefined;
  }

  return moment(p.TimeStart).add(hh, 'hours')
                            .add(mm, 'minutes')
                            .add(mm, 'seconds');
}

export default InvoicesStore;

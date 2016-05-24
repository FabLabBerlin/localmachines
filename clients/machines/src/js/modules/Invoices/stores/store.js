var $ = require('jquery');
var _ = require('lodash');
var actionTypes = require('../actionTypes');
var moment = require('moment');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;

const initialState = toImmutable({
  MonthlySums: {
  invoices: {},
  selected: {
      month: moment().month(),
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
    this.on(actionTypes.FETCH_MONTHLY_SUMMARIES, fetchMonthlySums);
    this.on(actionTypes.SELECT_USER_ID, selectUserId);
    this.on(actionTypes.SET_SELECTED_MONTH, setSelectedMonth);
    this.on(actionTypes.SET_INVOICE, setInvoice);
    this.on(actionTypes.SET_INVOICE_STATUSES, setInvoiceStatuses);
    this.on(actionTypes.SET_USER_MEMBERSHIPS, setUserMemberships);
  }

});

function editPurchase(state, id) {
  return state.set('editPurchaseId', id);
}

function fetchMonthlySums(state, { month, year, summaries }) {
  return state.setIn(['MonthlySums', year, month], toImmutable(summaries));
}

function selectUserId(state, userId) {
  return state.setIn(['MonthlySums', 'selected', 'userId'], userId);
}

function setSelectedMonth(state, { month, year }) {
  return state.setIn(['MonthlySums', 'selected', 'month'], month)
              .setIn(['MonthlySums', 'selected', 'year'], year);
}

function setInvoice(state, { month, year, userId, invoice }) {
  return state.setIn(['invoices', year, month, userId], invoice);
}

function setInvoiceStatuses(state, { month, year, userId, invoiceStatuses }) {
  return state.setIn(['invoiceStatuses', year, month, userId], invoiceStatuses);
}

function setUserMemberships(state, { userId, userMemberships }) {
  return state.setIn(['userMemberships', userId], userMemberships);
}

export default InvoicesStore;

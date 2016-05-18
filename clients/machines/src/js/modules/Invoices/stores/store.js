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
    this.on(actionTypes.FETCH_MONTHLY_SUMMARIES, fetchMonthlySums);
    this.on(actionTypes.SELECT_USER_ID, selectUserId);
    this.on(actionTypes.SET_SELECTED_MONTH, setSelectedMonth);
    this.on(actionTypes.SET_INVOICE, setInvoice);
  }

});

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

export default InvoicesStore;

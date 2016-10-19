var _ = require('lodash');
var actionTypes = require('../actionTypes');
var getters = require('../getters');
var moment = require('moment');
var Nuclear = require('nuclear-js');
var reactor = require('../../../reactor');
var toImmutable = Nuclear.toImmutable;
var {formatDuration, toQuantity} = require('../../../components/UserProfile/helpers');

const initialState = toImmutable({
  MonthlySums: {
    invoices: {},
    selected: {
        month: moment().month() + 1,
        year: moment().year(),
        userId: undefined
    },
    sorting: {
      column: 'Name',
      asc: true
    }
  },
  checkStatus: 'all',
  invoices: {
    detailedInvoices: {}
  },
  showInactiveUsers: false
});


var InvoicesStore = new Nuclear.Store({

  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.CHECK, check);
    this.on(actionTypes.CHECK_ALL, checkAll);
    this.on(actionTypes.CHECK_SET_STATUS, checkSetStatus);
    this.on(actionTypes.EDIT_PURCHASE, editPurchase);
    this.on(actionTypes.EDIT_PURCHASE_DURATION, editPurchaseDuration);
    this.on(actionTypes.EDIT_PURCHASE_FIELD, editPurchaseField);
    this.on(actionTypes.FETCH_MONTHLY_SUMMARIES, fetchMonthlySums);
    this.on(actionTypes.SET_SELECTED_MONTH, setSelectedMonth);
    this.on(actionTypes.SET_INVOICE, setInvoice);
    this.on(actionTypes.SET_INVOICE_STATUSES, setInvoiceStatuses);
    this.on(actionTypes.SET_SHOW_INACTIVE_USERS, setShowInactiveUsers);
    this.on(actionTypes.SET_USER_MEMBERSHIPS, setUserMemberships);
    this.on(actionTypes.SORT_BY, sortBy);
  }

});

function check(state, id) {
  const month = state.getIn(['MonthlySums', 'selected', 'month']);
  const year = state.getIn(['MonthlySums', 'selected', 'year']);
  const invoices = state.getIn(['MonthlySums', year, month]).map((inv) => {
    if (inv.get('Id') === id) {
      return inv.set('checked', !inv.get('checked'));
    } else {
      return inv;
    }
  });

  return state.setIn(['MonthlySums', year, month], invoices);
}

function checkAll(state) {
  return updateChecks(state, true);
}

function checkSetStatus(state, status) {
  return updateChecks(state.set('checkStatus', status), false);
}

function editPurchase(state, id) {
  return state.set('editPurchaseId', id);
}

function editPurchaseField(state, {field, value, invoiceId}) {
  var purchaseId = state.get('editPurchaseId');

  var keyPath = [
    'invoices',
    'detailedInvoices',
    invoiceId
  ];

  var iv = state.getIn(keyPath).toJS();
  iv.Purchases = _.map(iv.Purchases, (p) => {
    if (p.Id === purchaseId) {
      if (!p.edited) {
        p.edited = {};
      }

      if (field === 'Type' && p.PriceUnit !== 'gram') {
        const durationString = formatDuration(p);
        console.log('p.Machine.PriceUnit=', p.Machine.PriceUnit);
        switch (value) {
        case 'activation':
          p.PriceUnit = 'minute';
          break;
        case 'reservation':
          p.PriceUnit = '30 minutes';
          break;
        }

        p.Quantity = toQuantity(p, durationString);
      }

      p[field] = value;
    }
    return p;
  });

  return state.setIn(keyPath, toImmutable(iv));
}

function editPurchaseDuration(state, {duration, invoiceId}) {
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

      var quantity = toQuantity(p, duration);
      if (quantity) {
        p.editValid = true;
        p.Quantity = quantity;
      }
    }
    return p;
  });

  return state.setIn(keyPath, toImmutable(iv));
}

function fetchMonthlySums(state, { month, year, summaries }) {
  return state.setIn(['MonthlySums', year, month], toImmutable(summaries));
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

function setShowInactiveUsers(state, { show }) {
  return state.set('showInactiveUsers', show);
}

function setUserMemberships(state, { userId, userMemberships }) {
  return state.setIn(['userMemberships', userId], userMemberships);
}

function sortBy(state, { column, asc }) {
  return state.setIn(['MonthlySums', 'sorting'], toImmutable({
    column: column,
    asc: asc
  }));
}

// Private:

function updateChecks(state, toggle) {
  const checkedAll = reactor.evaluateToJS(getters.getCheckedAll);
  const month = state.getIn(['MonthlySums', 'selected', 'month']);
  const year = state.getIn(['MonthlySums', 'selected', 'year']);
  const invoices = state.getIn(['MonthlySums', year, month])
  .map(inv => {
    var checked = toggle ? !checkedAll : checkedAll;
    switch (state.get('checkStatus')) {
    case 'all':
      break;
    case inv.get('Status'):
      break;
    default:
      checked = false;
      break;
    }
    return inv.set('checked', checked);
  });

  return state.setIn(['MonthlySums', year, month], invoices);  
}

export default InvoicesStore;

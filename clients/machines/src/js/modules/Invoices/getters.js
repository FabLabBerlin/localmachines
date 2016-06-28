var moment = require('moment');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;


const getEditPurchaseId = [
  ['invoicesStore'],
  (invoicesStore) => {
    return invoicesStore.get('editPurchaseId');
  }
];

const getMonthlySums = [
  ['invoicesStore'],
  (invoicesStore) => {
    return invoicesStore.get('MonthlySums');
  }
];

const getInvoice = [
	['invoicesStore'],
  (invoicesStore) => {
    const month = invoicesStore.getIn(['MonthlySums', 'selected', 'month']);
    const year = invoicesStore.getIn(['MonthlySums', 'selected', 'year']);
    const invoiceId = invoicesStore.getIn(['MonthlySums', 'selected', 'invoiceId']);
    const invoice = invoicesStore.getIn(['invoices', 'detailedInvoices', invoiceId]);

    return invoice;
  }
];

const getThisMonthInvoices = [
  getMonthlySums,
  (monthlySums) => {
    const month = monthlySums
                      .get('selected').get('month');
    const year = monthlySums
                      .get('selected').get('year');
    return monthlySums.getIn([year, month]);
  }
];

const getInvoiceActions = [
  getEditPurchaseId,
  getInvoice,
  (editPurchaseId, invoice) => {
    var as = {};

    const m = moment().month() + 1;
    const y = moment().year();

    const isPastMonth = invoice &&
     (y > invoice.get('Year') || m > invoice.get('Month'));
    const isPositive = invoice && invoice.get('Total') >= 0.01;

    if (invoice && y >= invoice.get('Year')) {
      switch (invoice.get('Status')) {
      case 'draft':
        console.log('inv=', invoice.toJS());
        as.Cancel = false;
        as.Freeze = true;
        as.PushDraft = true;
        as.Save = !!editPurchaseId;
        as.Send = false;
        break;
      case 'outgoing':
        if (invoice.get('Canceled')) {
          as.Cancel = false;
          as.Freeze = false;
          as.PushDraft = false;
          as.Save = false;
          as.Send = true;
        } else {
          as.Cancel = true;
          as.Freeze = false;
          as.PushDraft = false;
          as.Save = !!editPurchaseId;
          as.Send = true;
        }
        break;
      default:
        console.error('Unhandled status', invoice.get('Status'));
      }
    }

    if (!invoice || !invoice.getIn(['User', 'ClientId'])) {
      as.Cancel = false;
      as.Freeze = false;
      as.PushDraft = false;
      as.Send = false;
    }

    as.Freeze = as.Freeze && isPastMonth && isPositive;
    as.PushDraft = as.PushDraft && isPositive;

    return toImmutable(as);
  }
];

const getInvoiceStatuses = [
  ['invoicesStore'],
  (invoicesStore) => {
    const month = invoicesStore.getIn(['MonthlySums', 'selected', 'month']);
    const year = invoicesStore.getIn(['MonthlySums', 'selected', 'year']);
    const userId = invoicesStore.getIn(['MonthlySums', 'selected', 'userId']);
    const invoiceStatuses = invoicesStore.getIn(['invoiceStatuses', year, month, userId]);

    return invoiceStatuses;
  }
];

const getUserMemberships = [
  ['invoicesStore'],
  (invoicesStore) => {
    const userId = invoicesStore.getIn(['MonthlySums', 'selected', 'userId']);

    return invoicesStore.getIn(['userMemberships', userId]);
  }
];

const getCheckStatus = [
  ['invoicesStore'],
  (invoicesStore) => {
    return invoicesStore.get('checkStatus');
  }
];

const getCheckedAll = [
  getThisMonthInvoices,
  getCheckStatus,
  (invoices, checkStatus) => {
    if (invoices) {
      var allWithStatus = 0;
      const checked = invoices.reduce((n, inv) => {
        if (checkStatus === 'all' || checkStatus === inv.get('Status')) {
          allWithStatus++;
        }
        return n + (inv.get('checked') ? 1 : 0);
      }, 0);
      console.log('allWithStatus=', allWithStatus);
      console.log('checked=', checked);

      return allWithStatus > 0 && allWithStatus === checked;
    } else {
      return false;
    }
  }
];

export default {
  getCheckedAll,
  getCheckStatus,
  getEditPurchaseId,
  getInvoice,
  getInvoiceActions,
  getInvoiceStatuses,
  getMonthlySums,
  getUserMemberships
};

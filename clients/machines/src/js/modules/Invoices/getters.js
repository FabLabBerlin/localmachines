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

const getInvoiceActions = [
  getEditPurchaseId,
  getInvoice,
  (editPurchaseId, invoice) => {
    var as = {};
    var m = moment().month() + 1;
    var y = moment().year();

    if (invoice && y >= invoice.get('Year')) {
      switch (invoice.get('Status')) {
      case 'draft':
        console.log('inv=', invoice.toJS());
        as.Cancel = false;
        as.Freeze = (y > invoice.get('Year') || m > invoice.get('Month')) &&
          invoice.get('Total') >= 0.01;
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

export default {
  getEditPurchaseId,
  getInvoice,
  getInvoiceActions,
  getInvoiceStatuses,
  getMonthlySums,
  getUserMemberships
};

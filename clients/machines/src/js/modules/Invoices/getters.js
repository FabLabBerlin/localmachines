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
  getInvoiceStatuses,
  getMonthlySums,
  getUserMemberships
};

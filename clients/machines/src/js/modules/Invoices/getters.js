var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;


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
    const userId = invoicesStore.getIn(['MonthlySums', 'selected', 'userId']);
    const invoice = invoicesStore.getIn(['invoices', year, month, userId]);

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
  getInvoice,
  getInvoiceStatuses,
  getMonthlySums,
  getUserMemberships
};

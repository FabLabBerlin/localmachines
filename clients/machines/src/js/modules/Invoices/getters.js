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

export default {
  getInvoice,
  getMonthlySums
};

var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;


const getMonthlySums = [
  ['invoicesStore'],
  (invoicesStore) => {
    return invoicesStore.get('MonthlySums');
  }
];

export default {
  getMonthlySums
};

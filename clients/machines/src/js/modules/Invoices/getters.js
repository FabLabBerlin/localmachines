var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;


const getMonthlySummaries = [
  ['invoicesStore'],
  (invoicesStore) => {
    return invoicesStore.get('monthlySummaries');
  }
];

export default {
  getMonthlySummaries
};

var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;


const getAll = [
  ['categoriesStore'],
  (categoriesStore) => {
  	if (categoriesStore) {
      return categoriesStore.push(toImmutable({
        Id: 0,
        ShortName: 'other',
        Name: 'Other'
      }));
    }
  }
];

export default {
  getAll
};

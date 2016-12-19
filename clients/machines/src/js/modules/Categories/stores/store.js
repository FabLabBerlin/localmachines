var actionTypes = require('../actionTypes');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;


const initialState = null;

var CategoriesStore = new Nuclear.Store({
  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.ADD_CATEGORY, addCategory);
    this.on(actionTypes.SET_CATEGORIES, setCategories);
    this.on(actionTypes.UPDATE_CATEGORY, updateCategory);
  }
});

function addCategory(state, category) {
  return state.push(toImmutable(category));
}

function setCategories(state, categories) {
  return toImmutable(categories);
}

function updateCategory(state, category) {
  return state.map(c => {
    if (c.get('Id') === category.get('Id')) {
      return category;
    } else {
      return c;
    }
  });
}

export default CategoriesStore;

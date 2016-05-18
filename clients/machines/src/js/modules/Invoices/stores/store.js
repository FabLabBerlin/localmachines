var $ = require('jquery');
var _ = require('lodash');
var actionTypes = require('../actionTypes');
var moment = require('moment');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;

const initialState = toImmutable({
  MonthlySums: {
  selected: {
      month: moment().month(),
      year: moment().year()
    }
  }
});


var InvoicesStore = new Nuclear.Store({

  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.FETCH_MONTHLY_SUMMARIES, fetchMonthlySums);
    this.on(actionTypes.SET_SELECTED_MONTH, setSelectedMonth);
  }

});

function fetchMonthlySums(state, { month, year, summaries }) {
  return state.setIn(['MonthlySums', year, month], toImmutable(summaries));
}

function setSelectedMonth(state, { month, year }) {
  return state.setIn(['MonthlySums', 'selected'],
               toImmutable({ month: month, year: year }));
}

export default InvoicesStore;

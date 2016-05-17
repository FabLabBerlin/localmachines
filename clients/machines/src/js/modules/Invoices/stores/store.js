var $ = require('jquery');
var _ = require('lodash');
var actionTypes = require('../actionTypes');
var moment = require('moment');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;

const initialState = toImmutable({
  monthlySummaries: {
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
    this.on(actionTypes.FETCH_MONTHLY_SUMMARIES, fetchMonthlySummaries);
    this.on(actionTypes.SET_SELECTED_MONTH, setSelectedMonth);
  }

});

function fetchMonthlySummaries(state, { month, year, summaries }) {
  return state.setIn(['monthlySummaries', year, month], toImmutable(summaries));
}

function setSelectedMonth(state, { month, year }) {
  return state.setIn(['monthlySummaries', 'selected'],
               toImmutable({ month: month, year: year }));
}

export default InvoicesStore;

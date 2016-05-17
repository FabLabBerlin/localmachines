var $ = require('jquery');
var _ = require('lodash');
var actionTypes = require('../actionTypes');
var moment = require('moment');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;

const initialState = toImmutable({
  monthlySummaries: {
  selected: {
      month: moment().month() + 1,
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
  }

});

function fetchMonthlySummaries(state, { month, year, summaries }) {
  return state.setIn(['monthlySummaries', year, month], toImmutable(summaries));
}

export default InvoicesStore;

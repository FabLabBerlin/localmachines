var actionTypes = require('../actionTypes');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;


const initialState = toImmutable({
  reservationRules: undefined
});

var ReservationRulesStore = new Nuclear.Store({
  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.SET_RESERVATION_RULES, setReservationRules);
  }
});

function setReservationRules(state, reservationRules) {
  return state.set('reservationRules', toImmutable(reservationRules));
}

export default ReservationRulesStore;

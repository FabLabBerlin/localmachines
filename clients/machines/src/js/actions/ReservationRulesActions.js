var actionTypes = require('../actionTypes');
var ApiActions = require('./ApiActions');
var reactor = require('../reactor');

var ReservationRulesActions = {
  load() {
    ApiActions.getCall('/api/reservation_rules', function(reservationRules) {
      reactor.dispatch(actionTypes.SET_RESERVATION_RULES, reservationRules);
    });
  }
};

export default ReservationRulesActions;

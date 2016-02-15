var actionTypes = require('../actionTypes');
var ApiActions = require('./ApiActions');
var reactor = require('../reactor');

var ReservationRulesActions = {
  load(locationId) {
    var url = '/api/reservation_rules?location=' + locationId;
    ApiActions.getCall(url, function(reservationRules) {
      reactor.dispatch(actionTypes.SET_RESERVATION_RULES, reservationRules);
    });
  }
};

export default ReservationRulesActions;

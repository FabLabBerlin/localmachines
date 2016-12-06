var $ = require('jquery');
var actionTypes = require('../actionTypes');
var ApiActions = require('./ApiActions');
var reactor = require('../reactor');

var ReservationRulesActions = {
  load(locationId) {
    $.ajax({
      url: '/api/reservation_rules?location=' + locationId
    })
    .done(reservationRules => {
      reactor.dispatch(actionTypes.SET_RESERVATION_RULES, reservationRules);
    });
  }
};

export default ReservationRulesActions;

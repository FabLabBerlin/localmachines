var $ = require('jquery');
import actionTypes from '../actionTypes';
import reactor from '../reactor';

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

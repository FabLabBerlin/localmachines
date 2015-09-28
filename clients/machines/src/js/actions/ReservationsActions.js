var ReservationActions = {

  createEmpty() {
    reactor.dispatch(actionTypes.CREATE_EMPTY);
  },

  createSetMachine({ mid }) {
    reactor.dispatch(actionTypes.CREATE_SET_MACHINE);
  },

  createSetDate({ date }) {
    reactor.dispatch(actionTypes.CREATE_SET_DATE);
  },

  createSetTimeFrom({ timeFrom }) {
    reactor.dispatch(actionTypes.CREATE_SET_TIME_FROM);
  },

  createSetTimeTo({ timeTo }) {
    reactor.dispatch(actionTypes.CREATE_SET_TIME_TO);
  }
};

export default ReservationActions;

var actionTypes = require('../actionTypes');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;

const initialState = toImmutable({
  allMemberships: undefined
});


var MembershipsStore = new Nuclear.Store({

  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.SET_ALL_MEMBERSHIPS, setAllMemberships);
  }

});

function setAllMemberships(state, allMemberships) {
  return state.set('allMemberships', toImmutable(allMemberships));
}

export default MembershipsStore;

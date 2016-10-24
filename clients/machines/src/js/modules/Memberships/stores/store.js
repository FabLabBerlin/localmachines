var actionTypes = require('../actionTypes');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;

const initialState = toImmutable({
  allMemberships: undefined,
  showArchived: false
});


var MembershipsStore = new Nuclear.Store({

  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.SET_ALL_MEMBERSHIPS, setAllMemberships);
    this.on(actionTypes.SET_SHOW_ARCHIVED_MEMBERSHIPS, setShowArchived);
  }

});

function setAllMemberships(state, allMemberships) {
  return state.set('allMemberships', toImmutable(allMemberships));
}

function setShowArchived(state, yes) {
  return state.set('showArchived', yes);
}

export default MembershipsStore;

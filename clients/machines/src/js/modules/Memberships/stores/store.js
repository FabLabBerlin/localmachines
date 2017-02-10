var _ = require('lodash');
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
    this.on(actionTypes.SET_MEMBERSHIP_CATEGORY, setMembershipCategory);
  }

});

function setAllMemberships(state, allMemberships) {
  return state.set('allMemberships', toImmutable(allMemberships));
}

function setShowArchived(state, yes) {
  return state.set('showArchived', yes);
}

function setMembershipCategory(state, {membershipId, categoryId, yes}) {
  return state.set('allMemberships',
    state.get('allMemberships').map(mb => {
      if (mb.get('Id') === membershipId) {
        const ids = JSON.parse(mb.get('AffectedCategories') || '[]');
        const newIds = yes
          ? _.union(ids, [categoryId])
          : _.difference(ids, [categoryId]);
        
        return mb.set('AffectedCategories', '[' + newIds.join(',') + ']');
      } else {
        return mb;
      }
    }));
}

export default MembershipsStore;

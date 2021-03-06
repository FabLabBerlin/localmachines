import _ from 'lodash';
import actionTypes from '../actionTypes';
import Nuclear from 'nuclear-js';
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
    this.on(actionTypes.SET_MEMBERSHIP_ARCHIVE, setMembershipArchive);
    this.on(actionTypes.SET_MEMBERSHIP_CATEGORY, setMembershipCategory);
    this.on(actionTypes.SET_MEMBERSHIP_FIELD, setMembershipField);
  }

});

function setAllMemberships(state, allMemberships) {
  return state.set('allMemberships', toImmutable(allMemberships));
}

function setShowArchived(state, yes) {
  return state.set('showArchived', yes);
}

function setMembershipArchive(state, {membershipId, yes}) {
  return state.update('allMemberships', mbs => mbs.map(mb => {
    if (mb.get('Id') === membershipId) {
      return mb.set('Archived', yes);
    } else {
      return mb;
    }
  }));
}

function setMembershipCategory(state, {membershipId, categoryId, yes}) {
  return state.update('allMemberships', mbs => mbs.map(mb => {
    if (mb.get('Id') === membershipId) {
      const ids = JSON.parse(mb.get('AffectedCategories') || '[]');
      const op = yes ? _.union : _.difference;
      const newIds = op(ids, [categoryId]);
      
      return mb.set('AffectedCategories', '[' + newIds.join(',') + ']');
    } else {
      return mb;
    }
  }));
}

function setMembershipField(state, {membershipId, key, value}) {
  return state.update('allMemberships', mbs => mbs.map(mb => {
    if (mb.get('Id') === membershipId) {
      return mb.set(key, value);
    } else {
      return mb;
    }
  }));
}

export default MembershipsStore;

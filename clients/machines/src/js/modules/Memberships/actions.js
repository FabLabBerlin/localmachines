var _ = require('lodash');
var $ = require('jquery');
var actionTypes = require('./actionTypes');
var getters = require('./getters');
var GlobalActions = require('../../actions/GlobalActions');
var reactor = require('../../reactor');
var toastr = require('../../toastr');


function fetch({locationId}) {
  $.ajax({
    url: '/api/memberships?location=' + locationId,
    dataType: 'json',
    type: 'GET',
    success(memberships) {
      reactor.dispatch(actionTypes.SET_ALL_MEMBERSHIPS, memberships);
    },
    error(xhr, status, err) {
      toastr.error('Error getting the memberships');
    }
  });
}

function setShowArchived(yes) {
  reactor.dispatch(actionTypes.SET_SHOW_ARCHIVED_MEMBERSHIPS, yes);
}

function setMembershipArchive(membershipId, yes) {
  reactor.dispatch(actionTypes.SET_MEMBERSHIP_ARCHIVE, {membershipId, yes});
  save(membershipId);
}

function setMembershipField(membershipId, key, value) {
  reactor.dispatch(actionTypes.SET_MEMBERSHIP_FIELD, {membershipId, key, value});
}

function setMembershipCategory(membershipId, categoryId, yes) {
  reactor.dispatch(actionTypes.SET_MEMBERSHIP_CATEGORY, {membershipId, categoryId, yes});
}

function save(membershipId) {
  const memberships = reactor.evaluateToJS(getters.getAllMemberships);
  const membership = _.find(memberships, mb => mb.Id === membershipId);

  GlobalActions.showGlobalLoader();

  $.ajax({
    url: '/api/memberships/' + membership.Id + '?location=' + membership.LocationId,
    dataType: 'json',
    type: 'PUT',
    contentType: 'application/json; charset=utf-8',
    data: JSON.stringify(membership)
  })
  .done(function() {
    toastr.success('Membership updated');
  })
  .fail(function() {
    toastr.error('Failed to update membership');
  })
  .always(GlobalActions.hideGlobalLoader);
}

export default {
  fetch,
  setShowArchived,
  setMembershipArchive,
  setMembershipCategory,
  setMembershipField,
  save
};

var $ = require('jquery');
var actionTypes = require('./actionTypes');
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

export default {
  fetch
};

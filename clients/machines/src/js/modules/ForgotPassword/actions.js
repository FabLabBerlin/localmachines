var $ = require('jquery');
var actionTypes = require('./actionTypes');
var reactor = require('../../reactor');
var toastr = require('../../toastr');

function emailReset(router, email) {
  $.ajax({
    method: 'POST',
    url: '/api/users/forgot_password',
    data: {
      email: email
    }
  })
  .success(function() {
    router.transitionTo('/forgot_password/email_sent');
  })
  .error(function() {
    toastr.error('An error occurred.  Please try again later.');
  });
}

export default {
  emailReset
};

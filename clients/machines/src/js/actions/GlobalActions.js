var $ = require('jquery');
var actionTypes = require('../actionTypes');
var reactor = require('../reactor');
var toastr = require('../toastr');


function showGlobalLoader() {
  reactor.dispatch(actionTypes.SET_LOADING);
}

function hideGlobalLoader() {
  reactor.dispatch(actionTypes.UNSET_LOADING);
}

function performSubscribeNewsletter(email) {
  showGlobalLoader();
  $.ajax({
    method: 'POST',
    url: '/api/newsletters/easylab_dev',
    data: {
      email: email
    },
    success: function() {
      hideGlobalLoader();
      toastr.info('Please check your E-Mails to confirm the subscription.');
    },
    error: function() {
      hideGlobalLoader();
      toastr.error('An error occurred, please try again later.');
    }
  });
}

export default {
  showGlobalLoader, hideGlobalLoader, performSubscribeNewsletter
};

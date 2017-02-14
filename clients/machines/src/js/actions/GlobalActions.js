import _ from 'lodash';
var $ = require('jquery');
import actionTypes from '../actionTypes';
import reactor from '../reactor';
import toastr from '../toastr';


if ($(window)) {
  // We're not inside a unit test...

  $(window).resize(_.debounce(() => {
    setWindowSize({
      width: window.innerWidth,
      height: window.innerHeight
    });
  }, 500));
}

function setWindowSize({width, height}) {
  reactor.dispatch(actionTypes.SET_WINDOW_SIZE, {width, height});
}

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

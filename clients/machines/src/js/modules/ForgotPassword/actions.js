var $ = require('jquery');
var actionTypes = require('./actionTypes');
var getters = require('./getters');
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

function getParameterByName(name, url) {
    if (!url) {
      url = window.location.href;
    }
    name = name.replace(/[\[\]]/g, '\\$&');
    var regex = new RegExp('[?&]' + name + '(=([^&#]*)|&|#|$)'),
        results = regex.exec(url);
    if (!results) {
      return null;
    }
    if (!results[2]) {
      return '';
    }
    return decodeURIComponent(results[2].replace(/\+/g, ' '));
}

function handleServerErrors(xhr) {
  const msg = xhr.responseText;

  console.log('msg:', msg);

  if (xhr.status === 401 && msg === 'Outdated key') {
    toastr.error('The key is too old.  Please try again and hurry up :)');
  } else if (xhr.status === 401 && msg === 'Wrong key') {
    toastr.error('The url seems wrong.  Please check your spam folder or generate a new key.');
  } else if (xhr.status === 401 && msg === 'Wrong phone') {
    toastr.error('The phone number does not seem correct.');
  } else {
    toastr.error('An error occurred.  Please try again later.');
  }
}

function submitPhone(router, phone) {
  const key = getParameterByName('key');

  $.ajax({
    method: 'POST',
    url: '/api/users/forgot_password/phone',
    data: {
      key: key,
      phone: phone
    }
  })
  .success(function() {
    reactor.dispatch(actionTypes.SET_KEY, key);
    reactor.dispatch(actionTypes.SET_PHONE, phone);
    router.transitionTo('/forgot_password/reset');
  })
  .error(function(xhr, status) {
    handleServerErrors(xhr);
  });
}

function submitPassword(router, password) {
  const key = reactor.evaluateToJS(getters.getKey);
  const phone = reactor.evaluateToJS(getters.getPhone);

  $.ajax({
    method: 'POST',
    url: '/api/users/forgot_password/reset',
    data: {
      key: key,
      phone: phone,
      password: password
    }
  })
  .success(function() {
    router.transitionTo('/forgot_password/done');
  })
  .error(function(xhr, status) {
    handleServerErrors(xhr);
  });
}

export default {
  emailReset,
  submitPhone,
  submitPassword
};

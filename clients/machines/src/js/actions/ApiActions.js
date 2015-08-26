var $ = require('jquery');
var actionTypes = require('../actionTypes');
var reactor = require('../reactor');
var toastr = require('../toastr');


/*
 * GET call to the API
 * Make GET call cutomisable
 */
function getCall(url, successFunction, toastrMessage = '', errorFunction = function() {}) {
  $.ajax({
    url: url,
    dataType: 'json',
    type: 'GET',
    cache: false,
    success: function(data) {
      successFunction(data);
    },
    error: function(xhr, status, err) {
      if (toastrMessage) {
        toastr.error(toastrMessage);
      }
      errorFunction();
      console.error(url, status, err);
    }
  });
}

/*
 * POST call to the API
 * Make POST call cutomisable
 */
function postCall(url, dataToSend, successFunction, toastrMessage = '', errorFunction = function() {}) {
  showGlobalLoader();
  $.ajax({
    url: url,
    dataType: 'json',
    type: 'POST',
    data: dataToSend,
    success: function(data) {
      hideGlobalLoader();
      successFunction(data);
    },
    error: function(xhr, status, err) {
      hideGlobalLoader();
      if (toastrMessage) {
        toastr.error(toastrMessage);
      }
      errorFunction();
      console.error(url, status, err);
    }
  });
}

function showGlobalLoader() {
  reactor.dispatch(actionTypes.SET_LOADING);
}

function hideGlobalLoader() {
  reactor.dispatch(actionTypes.UNSET_LOADING);
}

export default {
  getCall, postCall,
  showGlobalLoader, hideGlobalLoader
};

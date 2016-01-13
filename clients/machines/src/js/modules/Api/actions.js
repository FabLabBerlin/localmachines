var $ = require('jquery');
var GlobalActions = require('../Global/actions');
var reactor = require('../../reactor');
var toastr = require('../../toastr');


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
    success(data) {
      successFunction(data);
    },
    error(xhr, status, err) {
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
  GlobalActions.showGlobalLoader();
  $.ajax({
    url: url,
    dataType: 'json',
    type: 'POST',
    data: dataToSend,
    success(data) {
      GlobalActions.hideGlobalLoader();
      successFunction(data);
    },
    error(xhr, status, err) {
      GlobalActions.hideGlobalLoader();
      if (toastrMessage) {
        toastr.error(toastrMessage);
      }
      errorFunction();
      console.error(url, status, err);
    }
  });
}

export default {
  getCall, postCall
};

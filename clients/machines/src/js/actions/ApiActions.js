var $ = require('jquery');
var actionTypes = require('../actionTypes');
var GlobalActions = require('./GlobalActions');
var reactor = require('../reactor');
var toastr = require('../toastr');


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
  postCall
};

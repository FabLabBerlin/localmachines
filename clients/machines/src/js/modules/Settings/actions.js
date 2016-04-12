var $ = require('jquery');
var actionTypes = require('./actionTypes');
var Cookies = require('js-cookie');
var reactor = require('../../reactor');
var toastr = require('../../toastr');

var SettingsActions = {

  loadSettings({locationId}) {
    $.ajax({
      url: '/api/settings/vat_percent?location=' + locationId,
      dataType: 'json',
      success(vatPercent) {
        reactor.dispatch(actionTypes.SET_VAT_PERCENT, vatPercent);
      },
      error(xhr, status, err) {
        toastr.error('Error loading settings');
      }
    });
  }

};





export default SettingsActions;

var _ = require('lodash');
var $ = require('jquery');
var actionTypes = require('./actionTypes');
var Cookies = require('js-cookie');
var LocationGetters = require('../../modules/Location/getters');
var reactor = require('../../reactor');
var SettingsGetters = require('../../modules/Settings/getters');
var toastr = require('../../toastr');


var SettingsActions = {

  loadAdminSettings({locationId}) {
    $.ajax({
      url: '/api/settings?location=' + locationId,
      dataType: 'json',
      success(adminSettings) {
        var h = {
          Currency: {},
          TermsUrl: {},
          VAT: {}
        };
        _.each(adminSettings, function(setting) {
          h[setting.Name] = setting;
        });
        reactor.dispatch(actionTypes.SET_ADMIN_SETTINGS, h);
      },
      error(xhr, status, err) {
        toastr.error('Error loading admin settings');
      }
    });
  },

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
  },

  saveAdminSettings(adminSettings) {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    const data = JSON.stringify(_.map(adminSettings, function(setting, name) {
      return _.extend({
        Name: name
      }, setting);
    }));
    console.log('data => ', data);
    $.ajax({
      method: 'POST',
      url: '/api/settings?location=' + locationId,
      dataType: 'json',
      contentType: 'application/json',
      data: data
    })
    .success(() => {
      toastr.info('Successfully updated settings.');
      SettingsActions.loadSettings({locationId});
    })
    .error(() => {
      toastr.error('Error saving settings.');
    });
  }

};





export default SettingsActions;

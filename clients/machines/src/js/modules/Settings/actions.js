var _ = require('lodash');
var $ = require('jquery');
var actionTypes = require('./actionTypes');
var Cookies = require('js-cookie');
var LocationGetters = require('../../modules/Location/getters');
var reactor = require('../../reactor');
var SettingsGetters = require('../../modules/Settings/getters');
var toastr = require('../../toastr');


var SettingsActions = {

  loadSettings({locationId}) {
    $.ajax({
      url: '/api/settings?location=' + locationId,
      dataType: 'json',
      success(adminSettings) {
        var h = {
          Currency: {},
          TermsUrl: {},
          VAT: {},
          FastbillTemplateId: {}
        };
        _.each(adminSettings, function(setting) {
          h[setting.Name] = setting;
        });
        reactor.dispatch(actionTypes.SET_SETTINGS, h);
      },
      error(xhr, status, err) {
        toastr.error('Error loading admin settings');
      }
    });
  },

  loadFastbillTemplates({locationId}) {
    $.ajax({
      url: '/api/settings/fastbill_templates?location=' + locationId,
      success(fastbillTemplates) {
        reactor.dispatch(actionTypes.SET_FASTBILL_TEMPLATES, fastbillTemplates);
      },
      error(xhr, status, err) {
        toastr.error('Error loading fastbill templates');
      }
    });
  },

  saveAdminSettings(adminSettings) {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    console.log('dat=a', data);
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

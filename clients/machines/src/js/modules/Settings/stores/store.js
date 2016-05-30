var $ = require('jquery');
var actionTypes = require('../actionTypes');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;


const initialState = toImmutable({});


var SettingsStore = new Nuclear.Store({
  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.SET_ADMIN_SETTINGS, setAdminSettings);
    this.on(actionTypes.SET_VAT_PERCENT, setVatPercent);
  }
});

function setAdminSettings(state, adminSettings) {
  return state.set('adminSettings', adminSettings);
}

function setVatPercent(state, vatPercent) {
  return state.set('VatPercent', vatPercent);
}

export default SettingsStore;

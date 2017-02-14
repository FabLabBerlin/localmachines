var $ = require('jquery');
import actionTypes from '../actionTypes';
import Nuclear from 'nuclear-js';
var toImmutable = Nuclear.toImmutable;


const initialState = toImmutable({});


var SettingsStore = new Nuclear.Store({
  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.SET_SETTINGS, setSettings);
    this.on(actionTypes.SET_FASTBILL_TEMPLATES, setFastbillTemplates);
  }
});

function setSettings(state, settings) {
  return state.set('settings', toImmutable(settings));
}

function setFastbillTemplates(state, fastbillTemplates) {
  return state.set('fastbillTemplates', toImmutable(fastbillTemplates));
}

export default SettingsStore;

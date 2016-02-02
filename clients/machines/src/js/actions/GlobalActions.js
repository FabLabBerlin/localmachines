var $ = require('jquery');
var actionTypes = require('../actionTypes');

function showGlobalLoader() {
  reactor.dispatch(actionTypes.SET_LOADING);
}

function hideGlobalLoader() {
  reactor.dispatch(actionTypes.UNSET_LOADING);
}

export default {
  showGlobalLoader, hideGlobalLoader
};

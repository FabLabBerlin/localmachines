var actionTypes = require('../actionTypes');
var reactor = require('../reactor');


function showGlobalLoader() {
  reactor.dispatch(actionTypes.SET_LOADING);
}

function hideGlobalLoader() {
  reactor.dispatch(actionTypes.UNSET_LOADING);
}

export default {
  showGlobalLoader, hideGlobalLoader
};

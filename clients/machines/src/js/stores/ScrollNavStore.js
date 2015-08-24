var $ = require('jquery');
var actionTypes = require('../actionTypes');
var Nuclear = require('nuclear-js');
var toImmutable = Nuclear.toImmutable;

const initialState = toImmutable({
  animationRunning: false,
  position: 0,
  scrollStep: 100,
  upEnabled: false,
  downEnabled: !!window.libnfc
});

var ScrollNavStore = new Nuclear.Store({
  getInitialState() {
    return initialState;
  },

  initialize() {
    this.on(actionTypes.SCROLL_UP, scrollUp);
    this.on(actionTypes.SCROLL_DOWN, scrollDown);
    this.on(actionTypes.SET_SCROLL_STEP, setScrollStep);
  }
});

function scrollUp(state) {
  var pos = Math.max(0, state.get('position') - state.get('scrollStep'));
  return update(state.set('position', pos));
}

function scrollDown(state) {
  var pos = state.get('position') + state.get('scrollStep');
  if (pos + $(window).height() >= $('html,body').height()) {
    pos = $('html,body').height() - $(window).height();
  }
  return update(state.set('position', pos));
}

function update(state) {
  return state.set('upEnabled', state.get('position') > 0)
              .set('downEnabled', state.get('position') + $(window).height() < $('html,body').height());
}

function setScrollStep(state, { scrollStep }) {
  return state.set('scrollStep', scrollStep);
}

export default ScrollNavStore;

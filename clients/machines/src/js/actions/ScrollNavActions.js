var actionTypes = require('../actionTypes');
var getters = require('../getters');
var reactor = require('../reactor');


var ScrollNavActions = {
  scrollUp() {
    reactor.dispatch(actionTypes.SCROLL_UP);
    scrollAnimate();
  },

  scrollDown() {
    reactor.dispatch(actionTypes.SCROLL_DOWN);
    scrollAnimate();
  },

  setScrollStep(scrollStep) {
    reactor.dispatch(actionTypes.SET_SCROLL_STEP, { scrollStep });
  }
};

function scrollAnimate() {
  $('html,body').animate({
    scrollTop: reactor.evaluateToJS(getters.getScrollPosition)
  }, 'easeOutExpo');
}

export default ScrollNavActions;

var Login = require('../modules/Login');
var React = require('react');
var reactor = require('../reactor');
var ScrollNav = require('../modules/ScrollNav');

const TOP = 'top';
const BOTTOM = 'bottom';

var Button = React.createClass({
  handleClick() {
    Login.actions.keepAlive();
    switch (this.props.topOrBottom) {
    case TOP:
      ScrollNav.actions.scrollUp();
      break;
    case BOTTOM:
      ScrollNav.actions.scrollDown();
      break;
    }
  },

  render() {
    var faClassName = 'fa fa-arrow-circle-';
    if (this.props.topOrBottom === TOP) {
      faClassName += 'up';
    } else {
      faClassName += 'down';
    }
    return (
      <div className={'scroll-nav scroll-nav-' + this.props.topOrBottom}
           onClick={this.handleClick}>
        <i className={faClassName}/>
      </div>
    );
  }
});

export default React.createClass({
  mixins: [ reactor.ReactMixin ],

  componentDidMount() {
    window.setTimeout(function() {
      var scrollStep = $(window).height() / 2;
      ScrollNav.actions.setScrollStep(scrollStep);
    }, 200);
  },

  getDataBindings() {
    return {
      scrollUpEnabled: ScrollNav.getters.getScrollUpEnabled,
      scrollDownEnabled: ScrollNav.getters.getScrollDownEnabled
    };
  },

  render() {
    return (
      <div>
        {this.state.scrollUpEnabled ? <Button topOrBottom={TOP}/> : ''}
        {this.state.scrollDownEnabled ? <Button topOrBottom={BOTTOM}/> : ''}
      </div>
    );
  }
});


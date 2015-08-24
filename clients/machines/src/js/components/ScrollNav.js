var getters = require('../getters');
var React = require('react');
var reactor = require('../reactor');
var ScrollNavActions = require('../actions/ScrollNavActions');

const TOP = 'top';
const BOTTOM = 'bottom';

var Button = React.createClass({
  handleClick() {
    switch (this.props.topOrBottom) {
    case TOP:
      ScrollNavActions.scrollUp();
      break;
    case BOTTOM:
      ScrollNavActions.scrollDown();
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

var ScrollNav = React.createClass({
  mixins: [ reactor.ReactMixin ],

  componentDidMount() {
    ScrollNavActions.setScrollStep($(window).height() / 2);
  },

  getDataBindings() {
    return {
      scrollUpEnabled: getters.getScrollUpEnabled,
      scrollDownEnabled: getters.getScrollDownEnabled
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

export default ScrollNav;

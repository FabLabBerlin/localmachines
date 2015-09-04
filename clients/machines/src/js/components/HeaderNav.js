var LoginActions = require('../actions/LoginActions');
var {Navigation} = require('react-router');
var React = require('react');


var Button = React.createClass({
  render() {
    var className;
    if (this.props.href === (window.location.pathname + window.location.hash)) {
      className = 'header-nav-active';
    }
    return (
      <a id={this.props.id} className={className} href={this.props.href}>
        {this.props.children}
      </a>
    );
  }
});


var HeaderNav = React.createClass({
  mixins: [ Navigation ],

  handleClick() {
    LoginActions.logout(this.context.router);
  },

  render() {
    var buttons = [];
    if (!window.libnfc) {
      buttons.push(<Button id="header-nav-machines" href="/machines/#/machine">Machines</Button>);
      buttons.push(<Button id="header-nav-profile" href="/machines/#/profile">Profile</Button>);
      buttons.push(<Button id="header-nav-spendings" href="/machines/#/spendings">Spendings</Button>);
    }
    return (
      <div className="header-nav">
        {buttons}
        <button
          className="btn btn-danger btn-logout pull-right"
          onClick={this.handleClick}>
          <i className="fa fa-sign-out"></i>
        </button>
      </div>
    );
  }
});

export default HeaderNav;

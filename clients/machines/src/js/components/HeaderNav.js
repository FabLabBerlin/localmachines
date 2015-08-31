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
  render() {
    return (
      <div className="header-nav">
        <Button id="header-nav-machines" href="/machines/#/machine">Machines</Button>
        <Button id="header-nav-profile" href="/machines/#/profile">Profile</Button>
        <Button id="header-nav-spendings" href="/machines/#/spendings">Spendings</Button>
      </div>
    );
  }
});

export default HeaderNav;

var React = require('react');


var HeaderNav = React.createClass({
  render() {
    return (
      <div className="header-nav">
        <a href="/machines/#/machine">Machines</a>
        <a href="/machines/#/profile">Profile</a>
        <a href="/machines/#/spendings">Spendings</a>
      </div>
    );
  }
});

export default HeaderNav;

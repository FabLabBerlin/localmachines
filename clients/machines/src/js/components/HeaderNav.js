var React = require('react');


var HeaderNav = React.createClass({
  render() {
    return (
      <div className="header-nav">
        <a id="header-nav-machines" href="/machines/#/machine">Machines</a>
        <a id="header-nav-profile" href="/machines/#/profile">Profile</a>
        <a id="header-nav-spendings" href="/machines/#/spendings">Spendings</a>
      </div>
    );
  }
});

export default HeaderNav;

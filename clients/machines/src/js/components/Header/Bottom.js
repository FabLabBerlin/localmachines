var React = require('react');


var Item = React.createClass({
  render() {
    const isActive = this.props.href === (window.location.pathname + '#' + this.props.location.pathname) ||
      /* TODO: remove, just for ongoing works */
      this.props.href + 's' === (window.location.pathname + '#' + this.props.location.pathname);
    const activeClass = isActive ? 'active' : '';

    return (
      <div className="nav-item-container">
        <div className={'nav-item ' + activeClass + ' ' + this.props.className}>
          <a href={this.props.href}>
            <span>{this.props.label}</span>
            <img src={this.props.icon}/>
          </a>
        </div>
      </div>
    );
  }
});


var Bottom = React.createClass({
  render() {
    return (
      <div className="nav-bottom row">
        <Item className="nav-item-machines"
              href="/machines/#/machine"
              icon="/machines/assets/img/header_nav/machine.svg"
              label="Machines"
              location={this.props.location}/>
        <Item className="nav-item-reservations"
              href="/machines/#/reservations"
              icon="/machines/assets/img/header_nav/reservations.svg"
              label="Reservations"
              location={this.props.location}/>
        <Item className="nav-item-spendings"
              href="/machines/#/spendings"
              icon="/machines/assets/img/header_nav/spendings.svg"
              label="Spendings"
              location={this.props.location}/>
      </div>
    );
  }
});

export default Bottom;

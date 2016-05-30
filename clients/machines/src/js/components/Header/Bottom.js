var React = require('react');


var Item = React.createClass({
  render() {
    var activeClass = '';
    if (this.props.href === (window.location.pathname + window.location.hash)) {
      activeClass = 'active';
    }

    return (
      <div className="nav-item-container">
        <div className={'nav-item ' + activeClass}>
          <a href={this.props.href}>
            {this.props.label}
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
        <Item label={<span><i className="fa fa-plug"/> <span className="hidden-xs">Machines</span></span>}
              href="/machines/#/machine"/>
        <Item label={<span><i className="fa fa-calendar-check-o"/> <span className="hidden-xs">Reservations</span></span>}
              href="/machines/#/reservations"/>
        <Item label={<span><i className="fa fa-money"/> <span className="hidden-xs">Spendings</span></span>}
              href="/machines/#/spendings"/>
      </div>
    );
  }
});

export default Bottom;

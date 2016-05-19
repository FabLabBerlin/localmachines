var React = require('react');


var Item = React.createClass({
  render() {
    return (
      <div className="col-xs-4 text-center">
        <div className="nav-item">
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
      <div className="row">
        <Item label="Machines"
              href="/machines/#/machine"/>
        <Item label="Reservations"
              href="/machines/#/reservations"/>
        <Item label="Spendings"
              href="/machines/#/spendings"/>
      </div>
    );
  }
});

export default Bottom;

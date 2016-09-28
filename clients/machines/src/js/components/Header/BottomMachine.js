var Item = require('./Item');
var React = require('react');


var BottomMachine = React.createClass({
  render() {
    return (
      <div className="nav-bottom row">
        <Item className="nav-item-machines"
              href="/machines/#/machine"
              icon="/machines/assets/img/header_nav/machine.svg"
              label="Use"
              location={this.props.location}/>
        <Item className="nav-item-reservations"
              href="/machines/#/reservations"
              icon="/machines/assets/img/header_nav/reservations.svg"
              label="Reservation"
              location={this.props.location}/>
        <Item className="nav-item-spendings"
              href="/machines/#/spendings"
              icon="/machines/assets/img/header_nav/spendings.svg"
              label="Infos"
              location={this.props.location}/>
      </div>
    );
  }
});

export default BottomMachine;

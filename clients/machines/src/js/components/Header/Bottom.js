import Item from './Item';
import React from 'react';


var Bottom = React.createClass({
  render() {
    return (
      <div className="nav-bottom row">
        <Item className="nav-item-machines"
              cols={3}
              href="/machines/#/machines"
              icon="/machines/assets/img/header_nav/machine.svg"
              label="Machines"
              location={this.props.location}/>
        <Item className="nav-item-reservations"
              cols={3}
              href="/machines/#/reservations"
              icon="/machines/assets/img/header_nav/reservations.svg"
              label="Reservations"
              location={this.props.location}/>
        <Item className="nav-item-spendings"
              cols={3}
              href="/machines/#/spendings"
              icon="/machines/assets/img/header_nav/spendings.svg"
              label="Spendings"
              location={this.props.location}/>
      </div>
    );
  }
});

export default Bottom;

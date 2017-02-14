import _ from 'lodash';
import Item from './Item';
import Machines from '../../modules/Machines';
import React from 'react';
import reactor from '../../reactor';


var BottomMachine = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      machines: Machines.getters.getMachines
    };
  },

  render() {
    var m;

    if (this.state.machines) {
      m = this.state.machines.find((mm) => {
        return mm.get('Id') === this.props.machineId;
      });
    }

    if (!m || !_.isNumber(m.get('ReservationPriceHourly'))) {
      return (
        <div className="nav-bottom row">
          <Item className="nav-item-machines"
                cols={2}
                href={'/machines/#/machines/' + this.props.machineId}
                icon="/machines/assets/img/header_nav/machine.svg"
                label="Use"
                location={this.props.location}/>
          <Item className="nav-item-spendings"
                cols={2}
                href={'/machines/#/machines/' + this.props.machineId + '/infos'}
                icon="/machines/assets/img/header_nav/spendings.svg"
                label="Infos"
                location={this.props.location}/>
        </div>
      );
    } else {
      return (
        <div className="nav-bottom row">
          <Item className="nav-item-machines"
                cols={3}
                href={'/machines/#/machines/' + this.props.machineId}
                icon="/machines/assets/img/header_nav/machine.svg"
                label="Use"
                location={this.props.location}/>
          <Item className="nav-item-reservations"
                cols={3}
                href={'/machines/#/machines/' + this.props.machineId + '/reservations'}
                icon="/machines/assets/img/header_nav/reservations.svg"
                label="Reservation"
                location={this.props.location}/>
          <Item className="nav-item-spendings"
                cols={3}
                href={'/machines/#/machines/' + this.props.machineId + '/infos'}
                icon="/machines/assets/img/header_nav/spendings.svg"
                label="Infos"
                location={this.props.location}/>
        </div>
      );
    }
  }
});

export default BottomMachine;

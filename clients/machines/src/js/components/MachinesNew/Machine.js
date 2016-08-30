import {hashHistory} from 'react-router';
var getters = require('../../getters');
var React = require('react');
var reactor = require('../../reactor');
var Timer = require('../MachinePage/Machine/Timer');


const AVAILABLE = 'available';
const LOCKED = 'locked';
const MAINTENANCE = 'maintenance';
const OCCUPIED = 'occupied';
const RESERVED = 'reserved';
const UPCOMING_RESERVATION = 'upcoming-reservation';
const RUNNING = 'running';
const STAFF = 'staff';


var Machine = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      reservationsByMachineId: getters.getActiveReservationsByMachineId,
      upcomingReservationsByMachineId: getters.getUpcomingReservationsByMachineId,
      user: getters.getUser
    };
  },

  click() {
    hashHistory.push();
  },

  imgUrl() {
    if (this.props.machine.get('Image')) {
      return '/files/' + this.props.machine.get('Image');
    } else {
      return '/machines/img/img-machine-placeholder.svg';
    }
  },

  overlayText() {
    switch (this.status()) {
    case LOCKED:
      return 'Unlock ?';
    case MAINTENANCE:
      return 'Maintenance';
    case OCCUPIED:
      return 'Occupied';
    case RESERVED:
      return 'Reserved';
    case RUNNING:
      return (
        <div>
          <div>Running for</div>
          <Timer activation={this.props.machine.get('activation').toJS()}/>
        </div>
      );
    case UPCOMING_RESERVATION:
      return (
        <div>
          <div>Reserved in</div>
          <Timer activation={this.upcomingReservation().toJS()}/>
        </div>
      );
    }
  },

  render() {
    const m = this.props.machine;
    const style = {
      backgroundImage: 'url(' + this.imgUrl() + ')'
    };

    return (
      <a className={'ms-machine ' + this.statusClass()}
         href={'/machines/#/machines/' + m.get('Id')}>
        <div className="ms-machine-label">
          <div className="ms-machine-name">
            {m.get('Name')}
          </div>
          <div className="ms-machine-brand">
            {m.get('Brand')}
          </div>
        </div>
        <div className="ms-machine-icon" style={style}>
          <div className="ms-machine-overlay-container">
            <div className="ms-machine-overlay">
              {this.overlayText()}
            </div>
          </div>
        </div>
      </a>
    );
  },

  reservation() {
    const mid = this.props.machine.get('Id');

    if (this.state.reservationsByMachineId) {
      return this.state.reservationsByMachineId.toObject()[mid];
    }
  },

  upcomingReservation() {
    const mid = this.props.machine.get('Id');

    if (this.state.upcomingReservationsByMachineId) {
      return this.state.upcomingReservationsByMachineId.toObject()[mid];
    }
  },

  status() {
    const m = this.props.machine;
    const a = m.get('activation');
    const r = this.reservation();
    const upcoming = this.upcomingReservation();

    if (m.get('Locked')) {
      console.log('SOMETHING IS LOCKED!!!');
      return LOCKED;
    } else if (m.get('UnderMaintenance')) {
      return MAINTENANCE;
    } else {
      if (a) {
        if (a.get('UserId') === this.state.user.get('Id')) {
          return RUNNING;
        } else {
          return OCCUPIED;
        }
      } else if (r && !r.get('ReservationDisabled') && !r.get('Cancelled')) {
        return RESERVED;
      } else if (upcoming && !upcoming.get('ReservationDisabled') &&
                !upcoming.get('Cancelled')) {
        return UPCOMING_RESERVATION;
      } else {
        return AVAILABLE;
      }
    }
  },

  statusClass() {
    return 'ms-' + this.status();
  }
});

export default Machine;

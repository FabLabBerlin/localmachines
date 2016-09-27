var constants = require('./constants');
var React = require('react');
var ActivationTimer = require('../MachinePage/Machine/ActivationTimer');
var ReservationTimer = require('../MachinePage/Machine/ReservationTimer');


var MachineMixin = {
  imgUrl() {
    const m = this.machine();

    if (m.get('ImageSmall')) {
      return '/files/' + m.get('ImageSmall');
    } else {
      return '/machines/img/img-machine-placeholder.svg';
    }
  },

  overlayText() {
    switch (this.status()) {
    case constants.LOCKED:
      return 'Unlock ?';
    case constants.MAINTENANCE:
      return 'Maintenance';
    case constants.OCCUPIED:
      return 'Occupied';
    case constants.RESERVED:
      return 'Reserved';
    case constants.RUNNING:
      return (
        <div>
          <div className="ms-machine-overlay-start">Running for</div>
          <ActivationTimer activation={this.props.machine.get('activation').toJS()}/>
        </div>
      );
    case constants.UPCOMING_RESERVATION:
      return (
        <div>
          <div className="ms-machine-overlay-start">Reserved in</div>
          <ReservationTimer reservation={this.upcomingReservation().toJS()}/>
        </div>
      );
    }
  },

  status() {
    const m = this.props.machine;

    if (!m) {
      return undefined;
    }

    const a = m.get('activation');
    const r = this.reservation();
    const upcoming = this.upcomingReservation();

    if (m.get('Locked')) {
      return constants.LOCKED;
    } else if (m.get('UnderMaintenance')) {
      return constants.MAINTENANCE;
    } else {
      if (a) {
        if (a.get('UserId') === this.state.user.get('Id')) {
          return constants.RUNNING;
        } else {
          return constants.OCCUPIED;
        }
      } else if (r && !r.get('ReservationDisabled') && !r.get('Cancelled')) {
        return constants.RESERVED;
      } else if (upcoming && !upcoming.get('ReservationDisabled') &&
                !upcoming.get('Cancelled')) {
        return constants.UPCOMING_RESERVATION;
      } else {
        return constants.AVAILABLE;
      }
    }
  },

  statusClass() {
    return 'ms-' + this.status();
  }
};

export default MachineMixin;

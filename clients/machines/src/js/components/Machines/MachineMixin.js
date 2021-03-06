import constants from './constants';
import React from 'react';
import ActivationTimer from './ActivationTimer';
import ReservationTimer from './ReservationTimer';


var MachineMixin = {
  imgUrl(small) {
    const m = this.machine();
    const key = small ? 'ImageSmall' : 'Image';

    if (m.get(key)) {
      return '/files/' + m.get(key);
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
          <ActivationTimer activation={this.machine().get('activation').toJS()}/>
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

  status(ignoreUpcoming) {
    const m = this.machine();

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
      if (a && this.state.user) {
        if (a.get('UserId') === this.state.user.get('Id')) {
          return constants.RUNNING;
        } else {
          return constants.OCCUPIED;
        }
      } else if (r && !r.get('ReservationDisabled') && !r.get('Cancelled')) {
        return constants.RESERVED;
      } else if (upcoming && !upcoming.get('ReservationDisabled') &&
                !upcoming.get('Cancelled') &&
                !ignoreUpcoming) {
        return constants.UPCOMING_RESERVATION;
      } else {
        return constants.AVAILABLE;
      }
    }
  },

  reservation() {
    const mid = this.machine().get('Id');

    if (this.state.reservationsByMachineId) {
      return this.state.reservationsByMachineId.toObject()[mid];
    }
  },

  upcomingReservation() {
    const mid = this.machine().get('Id');

    if (this.state.upcomingReservationsByMachineId) {
      return this.state.upcomingReservationsByMachineId.toObject()[mid];
    }
  },

  statusClass() {
    return 'ms-' + this.status();
  }
};

export default MachineMixin;

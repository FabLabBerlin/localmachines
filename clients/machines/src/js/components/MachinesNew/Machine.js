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
const RUNNING = 'running';
const STAFF = 'staff';


var Machine = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      reservationsByMachineId: getters.getActiveReservationsByMachineId,
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
    }
  },

  render() {
    const m = this.props.machine;
    const style = {
      backgroundImage: 'url(' + this.imgUrl() + ')'
    };
    if (m.get('activation')) {
      console.log('machine with activation:', m.toJS());
    }

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
          <div className="ms-machine-overlay">
            {this.overlayText()}
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

  status() {
    const m = this.props.machine;
    const a = m.get('activation');
    const r = this.reservation();

    if (m.get('UnderMaintenance')) {
      return MAINTENANCE;
    } else {
      if (a) {
        console.log('machine ' + m.get('Name'));
        console.log('a=', a);
        console.log('a.UserId=', a.get('UserId'));
        console.log('userId=', this.state.user.get('Id'));
        console.log('this.state.user=', this.state.user);
        if (a.get('UserId') === this.state.user.get('Id')) {
          return RUNNING;
        } else {
          return OCCUPIED;
        }
      } else if (r && !r.get('ReservationDisabled') && !r.get('Cancelled')) {
        return RESERVED;
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

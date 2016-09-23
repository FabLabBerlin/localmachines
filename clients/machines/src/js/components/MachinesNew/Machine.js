var constants = require('./constants');
import {hashHistory} from 'react-router';
var getters = require('../../getters');
var MachineMixin = require('./MachineMixin');
var React = require('react');
var reactor = require('../../reactor');
var ActivationTimer = require('../MachinePage/Machine/ActivationTimer');
var ReservationTimer = require('../MachinePage/Machine/ReservationTimer');


var Machine = React.createClass({

  mixins: [ MachineMixin, reactor.ReactMixin ],

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

  machine() {
    return this.props.machine;
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
          <div className="ms-machine-overlay-background-container">
            <div className="ms-machine-overlay-background"/>
          </div>
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
  }
});

export default Machine;

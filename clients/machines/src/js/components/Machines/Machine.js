import constants from './constants';
import {hashHistory} from 'react-router';
import getters from '../../getters';
import MachineMixin from './MachineMixin';
import React from 'react';
import reactor from '../../reactor';


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
      backgroundImage: 'url(' + this.imgUrl(true) + ')'
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
  }
});

export default Machine;

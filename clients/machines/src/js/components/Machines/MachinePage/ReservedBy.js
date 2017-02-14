import getters from '../../../getters';
import LoaderLocal from '../../LoaderLocal';
import Machines from '../../../modules/Machines';
import moment from 'moment';
import React from 'react';
import reactor from '../../../reactor';

import helpers from '../../UserProfile/helpers';


var ReservedBy = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      machineUsers: Machines.getters.getMachineUsers,
      uid: getters.getUid
    };
  },

  render() {
    if (this.props.reservation) {
      if (this.state.machineUsers) {
        const uid = this.props.reservation.get('UserId');
        const users = this.state.machineUsers;
        const user = users.get(uid) || {};
        var timeEnd = helpers.timeEnd(this.props.reservation);
        if (timeEnd) {
          timeEnd = timeEnd.format('HH:mm');
        }

        return (
          <div className="m-indicator m-reserved-by">
            {uid === this.state.uid ? (
              <div>
                <div>My</div>
                <div>Reservation</div>
                <div>until {timeEnd}</div>
              </div>
            ) : (
              <div>
                <div>Reserved by</div>
                <div>{user.FirstName} {user.LastName}</div>
              </div>
            )}
          </div>
        );
      } else {
        return <LoaderLocal/>;
      }
    } else {
      return <div/>;
    }
  }
});

export default ReservedBy;

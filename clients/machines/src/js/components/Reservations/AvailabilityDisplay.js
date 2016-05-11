var getters = require('../../getters');
var Machines = require('../../modules/Machines');
var moment = require('moment');
var React = require('react');
var reactor = require('../../reactor');
var toastr = require('../../toastr');

const DELTA = 2.08333333;

function isNow(t) {
  var now = moment().unix();
  var nn = now - (now % 1800);
  var uu = t.unix() - (t.unix() % 1800);
  return nn === uu;
}


var Slot = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      userId: getters.getUid,
      machineUsers: Machines.getters.getMachineUsers
    };
  },

  handleClickOnReservedSlot(displayString) {
    toastr.info(displayString);
  },

  render() {
    var reserved = true;
    var reservedByUser = this.props.reservation && 
      this.props.reservation.get('UserId') === 
      this.state.userId;
    var users = this.state.machineUsers;

    var className = 'slot';
    var clickHandler = {};
    if (reserved) {
      clickHandler = this.handleClickOnReservedSlot;
      className += ' reserved';
      if (reservedByUser) {
        className += ' by-user';
      }
    }
    if (isNow(this.props.time) && !this.props.reservation) {
      className = 'slot now';
    }

    var title;
    var width = String(DELTA) + '%';
    if (this.props.reservation) {
      var start = moment(this.props.reservation.get('TimeStart'));
      var end = moment(this.props.reservation.get('TimeEnd'));
      width = String(DELTA * (end.unix() - start.unix()) / 1800) + '%';
      title = start.format('DD. MMM HH:mm') + ' - ' + end.format('HH:mm');
      var userId = this.props.reservation.get('UserId');
      var user = users.get(userId) || {};
      if (reservedByUser) {
        title += ' by You';
      } else {
        title += ' by ' + user.FirstName + ' ' + user.LastName;
      }
    }

    var style = {
      marginLeft: String(DELTA * this.props.position) + '%',
      width: width
    };

    return <div className={className}
                style={style}
                title={title}
                onClick={clickHandler.bind(this, title)}/>;
  }
});


var AvailabilityDisplay = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      reservations: getters.getReservations,
      slotAvailabilities48h: getters.getSlotAvailabilities48h,
      userId: getters.getUid
    };
  },

  render() {
    if (!this.state.reservations) {
      return <div/>;
    }

    var key = 1;

    var availabilities = this.state.slotAvailabilities48h.get(this.props.machineId);
    if (availabilities) {
      var today = availabilities.get('today');
      var tomorrow = availabilities.get('tomorrow');
      var todayStart = moment().hours(0);
      var todayEnd = todayStart.clone().add(1, 'day');
      var tomorrowStart = todayEnd.clone();
      var tomorrowEnd = tomorrowStart.clone().add(1, 'day');
      todayStart = todayStart.unix();
      todayEnd = todayEnd.unix();
      tomorrowStart = tomorrowStart.unix();
      tomorrowEnd = tomorrowEnd.unix();
      var indexNow = Math.round(2 * (moment().unix() - todayStart) / 3600);

      return (
        <div className="machine-reserv-preview">

          <div className="today">
            <div className="slots">
              {today.map(reservation => {
                var timeStart = moment(reservation.get('TimeStart'));
                var jj = Math.round(2 * (timeStart.unix() - todayStart) / (3600));
                return <Slot key={key++}
                             machineId={reservation.get('MachineId')}
                             position={jj}
                             reservation={reservation}
                             time={timeStart}/>;
              })}
              <Slot key={key++}
                    machineId={this.props.machineId}
                    position={indexNow}
                    time={moment()}/>
            </div> 

            <table className="label-bar">
              <tr>
                <td className="start-time">00:00</td>
                <td className="label">Availability Today</td>
                <td className="end-time">24:00</td>
              </tr>
            </table>

          </div>

          <div className="tomorrow hidden-xs">
            <div className="slots">
              {tomorrow.map(reservation => {
                var timeStart = moment(reservation.get('TimeStart'));
                var jj = Math.round(2 * (timeStart.unix() - tomorrowStart) / (3600));
                return <Slot key={key++}
                             machineId={reservation.get('MachineId')}
                             position={jj}
                             reservation={reservation}
                             time={timeStart}/>;
              })}
            </div>
            
            <table className="label-bar">
              <tr>
                <td className="start-time">00:00</td>
                <td className="label">Availability Tomorrow</td>
                <td className="end-time">24:00</td>
              </tr>
            </table>

          </div>

        </div>
      );
    } else {
      return <div/>;
    }
  }

});

export default AvailabilityDisplay;

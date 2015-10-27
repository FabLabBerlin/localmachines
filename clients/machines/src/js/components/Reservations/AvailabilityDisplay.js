var getters = require('../../getters');
var moment = require('moment');
var React = require('react');
var reactor = require('../../reactor');


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
      userId: getters.getUid
    };
  },

  render() {
    var reserved = true;
    var reservedByUser = this.props.reservation && this.props.reservation.get('UserId') === this.state.userId;

    var className = 'slot';
    if (reserved) {
      className += ' reserved';
      if (reservedByUser) {
        className += ' by-user';
      }
    }
    if (isNow(this.props.time)) {
      className += ' now';
    }

    var style = {
      marginLeft: String(2.08333333 * this.props.position) + '%'
    };

    return <div className={className}
                style={style}/>;
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

    var times = [];
    var i = 0;
    var n = 2 * 48;
    var reservation;

    for (var t = moment().hours(0); i < n; t.add(30, 'm'), i++) {
      times.push(t.clone());
    }

    var key = 1;

    var availabilities = this.state.slotAvailabilities48h.get(this.props.machineId);
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
          <div className="label">
            Today
          </div>
        </div>
        <div className="tomorrow">
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
          <div className="label">
            Tomorrow
          </div>
        </div>
      </div>
    );
  }

});

export default AvailabilityDisplay;

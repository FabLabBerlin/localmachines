var _ = require('lodash');
var getters = require('../../../../getters');
var LoaderLocal = require('../../../LoaderLocal');
var Machines = require('../../../../modules/Machines');
var moment = require('moment');
var React = require('react');
var reactor = require('../../../../reactor');


const MONDAY = 1;
const TUESDAY = 2;
const WEDNESDAY = 3;
const THURSDAY = 4;
const FRIDAY = 5;
const SATURDAY = 6;
const SUNDAY = 7;

function areSameDay(t1, t2) {
  return t1.format('YYYY-MM-DD') === t2.format('YYYY-MM-DD');
}

// E.g. toInt('13:00') = 2 * 13 = 26
function toInt(t) {
  if (t) {
    if ((typeof t === 'string' || t instanceof String) && t.length > 5) {
      t = moment(t);
    }
    if (moment.isMoment(t)) {
      t = t.format('HH:mm');
    }
    const tmp = t.split(':');
    var i = 2 * parseInt(tmp[0]);

    if (parseInt(tmp[1]) >= 30) {
      i += 1;
    }

    return i;
  }
}

function toTimeString(i) {
  const j = i % 2;
  const hh = (i - j) / 2;
  var mm = String(j * 30);

  if (mm.length === 1) {
    mm = '0' + mm;
  }

  return String(hh) + ':' + mm;
}

var Times = React.createClass({
  render() {
    var rows = [];

    for (var i = toInt(this.props.start); i < toInt(this.props.end); i++) {
      rows.push(
        <div key={i} className="r-time">
          {toTimeString(i)}
        </div>
      );
    }

    return (
      <div className="r-times">
        {rows}
      </div>
    );
  }
});


var DayHeader = React.createClass({
  render() {
    const isCurrentDay = areSameDay(this.props.day, moment()) ?
      'r-day-current' : '';

    return (
      <div className="r-day-header row">
        <div className="r-day-header-weekday col-xs-6">
          {this.props.day.format('ddd')}
        </div>
        <div className={'r-day-header-number col-xs-6 ' + isCurrentDay}>
          {this.props.day.format('D')}
        </div>
      </div>
    );
  }
});


var Slot = React.createClass({
  render() {
    return (
      <div className="r-slot"/>
    );
  }
});


var Event = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      machineUsers: Machines.getters.getMachineUsers
    };
  },

  render() {
    const uid = this.props.reservation.get('UserId');
    const users = this.state.machineUsers;
    const user = users.get(uid) || {};
    console.log('this.state.machineUsers=', this.state.machineUsers);
    const r = this.props.reservation;
    const i = toInt(r.get('TimeStart'));
    const j = toInt(r.get('TimeEnd'));
    const style = {
      height: (j - i) * 41
    };

    return (
      <div className="r-reservation" style={style}>
        {user.FirstName} {user.LastName}
      </div>
    );
  }
});


var Day = React.createClass({
  render() {
    if (!this.props.reservations) {
      return <LoaderLocal/>;
    }

    var rows = [];
    var key = 0;
    var r;
    var findReservation = (ii) => {
      return this.props.reservations.find(rr => {
        const j = toInt(rr.get('TimeStart'));

        return ii === j && this.props.machineId === rr.get('MachineId');
      });
    };

    for (var i = toInt(this.props.start); i < toInt(this.props.end); i++) {
      if (r) {
        const j = toInt(r.get('TimeStart'));
        const k = toInt(r.get('TimeEnd'));

        if (i < j || k <= i) {
          r = null;
        }
      }

      if (!r) {
        r = findReservation(i);
      }

      if (r && toInt(r.get('TimeStart')) === i) {
        rows.push(<Event key={++key} reservation={r}/>);
      } else if (!r) {
        rows.push(<Slot key={++key}/>);
      }
    }

    return (
      <div className="r-day">
        <DayHeader day={this.props.day}/>
        {rows}
      </div>
    );
  }
});


var Week = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      reservations: getters.getReservations,
      userId: getters.getUid
    };
  },

  render() {
    const start = '9:00';
    const end = '22:00';

    return (
      <div className="r-week">
        <Times start={start} end={end}/>
        {_.map(Array(7), (v, i) => {
          const day = this.props.startDay.clone().add(i, 'day');
          var reservations;

          if (this.state.reservations) {
            reservations = this.state.reservations.filter((r) => {
              return areSameDay(day, moment(r.get('TimeStart')));
            });
          }

          return <Day key={i}
                      day={day}
                      machineId={this.props.machineId}
                      reservations={reservations}
                      start={start}
                      end={end}/>;
        })}
      </div>
    );
  }
});


var Calendar = React.createClass({
  render() {
    const startDay = moment();

    while (startDay.weekday() !== MONDAY) {
      startDay.subtract(1, 'day');
    }

    return (
      <div id="r-calendar">
        <Week machineId={this.props.machineId} startDay={startDay}/>
      </div>
    );
  }
});

export default Calendar;

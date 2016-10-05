var _ = require('lodash');
var getters = require('../../../../getters');
var LoaderLocal = require('../../../LoaderLocal');
var MachineActions = require('../../../../actions/MachineActions');
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
      <div className="r-times visible-md-block visible-lg-block">
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


var Placeholder = React.createClass({
  render() {
    return (
      <div className="r-placeholder visible-md-block visible-lg-block"/>
    );
  }
});


var Slot = React.createClass({
  render() {
    return (
      <div className="r-slot visible-md-block visible-lg-block"/>
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
    const r = this.props.reservation;
    const i = toInt(r.get('TimeStart'));
    const j = toInt(r.get('TimeEnd'));
    const style = {
      height: (j - i) * 41
    };

    if (r.get('Name')) {
      return (
        <div className="r-reservation r-rule" style={style}>
          <div className="r-label">
            {r.get('Name')}
          </div>

          <div className="text-center visible-xs-block visible-sm-block">
            {r.get('TimeStart')} - {r.get('TimeEnd')}
          </div>
        </div>
      );
    } else {
      return (
        <div className="r-reservation" style={style}>
          <div className="r-label">
            {user.FirstName} {user.LastName}
          </div>

          <div className="text-center visible-xs-block visible-sm-block">
            {moment(r.get('TimeStart')).format('HH:mm')} - {moment(r.get('TimeEnd')).format('HH:mm')}
          </div>
        </div>
      );
    }
  }
});


var Day = React.createClass({
  render() {
    if (!this.props.reservations) {
      return <LoaderLocal/>;
    }

    const availableSlots = this.props.reservationRules
      .filter(rr => rr.get('Available'))
      .reduce((result, rr) => {

      for (var i = toInt(rr.get('TimeStart')); i < toInt(rr.get('TimeEnd')); i++) {
        result[i] = true;
      }

      return result;
    }, {});

    var rows = [];
    var key = 0;
    var res;
    var rule;
    var findReservation = (ii) => {
      return this.props.reservations.find(r => {
        const j = toInt(r.get('TimeStart'));

        return ii === j && this.props.machineId === r.get('MachineId');
      });
    };

    var findReservationRule = (ii) => {
      return this.props.reservationRules.find(rr => {
        if (rr.get('Unavailable')) {
          const j = toInt(rr.get('TimeStart'));

          return ii === j && (this.props.machineId === rr.get('MachineId')
                          || !rr.get('MachineId'));
        }
      });
    };

    for (var i = toInt(this.props.start); i < toInt(this.props.end); i++) {
      if (res) {
        const j = toInt(res.get('TimeStart'));
        const k = toInt(res.get('TimeEnd'));

        if (i < j || k <= i) {
          res = null;
        }
      }

      if (!res) {
        res = findReservation(i);
      }
      if (!res) {
        res = findReservationRule(i);
      }

      if (res && toInt(res.get('TimeStart')) === i) {
        rows.push(<Event key={++key} reservation={res}/>);
      } else if (!res) {
        if (availableSlots[i]) {
          rows.push(<Slot key={++key}/>);
        } else {
          rows.push(<Placeholder key={++key}/>);
        }
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
      reservationRules: getters.getReservationRules,
      userId: getters.getUid
    };
  },

  render() {
    var start = '14:00';
    var end = '17:00';

    if (!this.state.reservations || !this.state.reservationRules) {
      return <LoaderLocal/>;
    }

    const userIds = this.state.reservations
        .filter(r => r.get('MachineId') === this.props.machineId)
        .map(r => r.get('UserId'));

    MachineActions.fetchUserNames(userIds.toJS());

    this.state.reservationRules.forEach(rr => {
      if (rr.get('Available')) {
        if (toInt(rr.get('TimeStart')) < toInt(start)) {
          start = rr.get('TimeStart');
        }

        if (toInt(rr.get('TimeEnd')) > toInt(end)) {
          end = rr.get('TimeEnd');
        }
      }
    });

    return (
      <div className="r-week">
        <Times start={start} end={end}/>
        {_.map(Array(7), (v, i) => {
          const day = this.props.startDay.clone().add(i, 'day');
          const weekDayName = day.format('dddd');
          const reservations = this.state.reservations.filter((r) => {
            return areSameDay(day, moment(r.get('TimeStart')));
          });
          const reservationRules = this.state.reservationRules.filter((rr) => {
            const a = moment(rr.get('DateStart') + ' 00:00').unix();
            const b = moment(rr.get('DateEnd') + ' 00:00').add(1, 'day').unix();

            return a <= day.unix() && day.unix() <= b && rr.get(weekDayName);
          });

          return <Day key={i}
                      day={day}
                      machineId={this.props.machineId}
                      reservations={reservations}
                      reservationRules={reservationRules}
                      start={start}
                      end={end}/>;
        })}
      </div>
    );
  }
});


var Calendar = React.createClass({
  getInitialState() {
    const startDay = moment();

    while (startDay.weekday() !== MONDAY) {
      startDay.subtract(1, 'day');
    }

    return {
      startDay: startDay
    };
  },

  back() {
    this.setState({
      startDay: this.state.startDay.subtract(1, 'week')
    });
  },

  forward() {
    this.setState({
      startDay: this.state.startDay.add(1, 'week')
    });
  },

  clickCreate() {
    this.props.clickCreate();
  },

  render() {
    const endDay = this.state.startDay.clone().add(6, 'day');

    return (
      <div id="r-calendar">
        <div id="r-header">
          <button className="r-nav r-back" onClick={this.back}/>
          <button className="r-nav r-forward" onClick={this.forward}/>
          <span id="r-range">
            {this.state.startDay.format('MMM DD')} - {endDay.format('MMM DD, YYYY')}
          </span>
        </div>
        <div id="r-add-container" className="row">
          <div className="col-sm-6">
            <button id="r-add" onClick={this.clickCreate}/>
          </div>
          <div className="col-sm-6"/>
        </div>
        <div id="r-header-border"/>
        <Week machineId={this.props.machineId} startDay={this.state.startDay}/>
      </div>
    );
  }
});

export default Calendar;

import _ from 'lodash';
import getters from '../../getters';
import LocationGetters from '../../modules/Location/getters';
import moment from 'moment';
var { Month } = require('./helpers');
import React from 'react';
import reactor from '../../reactor';
import ReservationActions from '../../actions/ReservationActions';


var DayView = React.createClass({
  handleClick() {
    if (!this.props.header && !this.props.empty && !this.props.notAvailable) {
      var date = this.props.moment;
      ReservationActions.newReservation.setDate({ date });
    }
  },

  render() {
    if (this.props.header) {
      return (
        <div className="week-name">
          {this.props.day}
        </div>
      );
    } else if (this.props.empty) {
      return <div className="day empty"/>;
    } else {
      var className = 'day';
      if (this.props.notAvailable) {
        className += ' unavailable';
      } else {
        className += ' selectable';
      }
      if (this.props.moment.format('YYYY-MM-DD') === moment().format('YYYY-MM-DD')) {
        className += ' today';
      }
      if (this.props.selected) {
        className += ' selected';
      }
      return (
        <div className={className} onClick={this.handleClick}>
          {this.props.day}
        </div>
      );
    }
  }
});


/*
 * MonthView creates a view like:
 * M T W T F S S
 *       1 2 3 4
 *    . . .
 *
 * Note that MomentJS's week starts with Sunday = 0. Our
 * work week starts with Monday.
 * => WorkWeekDay = (MomentJSWeekDay + 6) % 7
 */
var MonthView = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      isAdmin: LocationGetters.getIsAdmin,
      user: getters.getUser
    };
  },

  render() {
    var selectedDate = reactor.evaluateToJS(getters.getNewReservation).date;
    var k = 0;
    var month = this.props.month;
    var days = [];
    var admin = this.state.isAdmin;
    for (var t = month.firstDay(); !t.isAfter(month.lastDay()); t = t.clone().add(1, 'day')) {
      days.push(t);
    }
    var weeks = [];
    while (days.length > 0) {
      if (weeks.length === 0) {
        var week = [];
        for (var j = 0; j % 7 < (days[0].day() + 6) % 7; j++) {
          week.push(<DayView key={k++} empty={true}/>);
        }
        weeks.push(week);
      }
      if (_.last(weeks).length === 7) {
        weeks.push([]);
      }
      var day = days.shift();
      _.last(weeks).push(
        <DayView key={k++}
                 day={day.date()}
                 moment={day}
                 selected={selectedDate && !selectedDate.diff(day, 'days')}
        />
      );
    }
    while (_.last(weeks).length < 7) {
      _.last(weeks).push(<DayView key={k++} empty={true}/>);
    }

    return (
      <div className="date-picker">
        <h4 className="month-name">{month.toString()}</h4>
        <div className="week">
          {_.map(['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'], function(heading, i) {
            return <DayView key={i} day={heading} header={true}/>;
          })}
        </div>
        {_.map(weeks, function(w, i) {
          return (
            <div key={i} className="week">
              {w}
            </div>
          );
        })}
        <div className="finish"/>
      </div>
    );
  }
});


var DatePicker = React.createClass({
  previous() {
    ReservationActions.newReservation.previousStep();
  },

  next() {
    ReservationActions.newReservation.nextStep();
  },

  render() {
    var currentMonth = Month.getCurrentMonth();
    var nextMonth = currentMonth.getNextMonth();
    return (
      <div className={this.props.className}>
        <h3>Select date</h3>
        <div id="date-picker">
          <MonthView month={currentMonth}/>
          <MonthView month={nextMonth}/>
        </div>
        <hr/>
        <div className="pull-right">
          <button className="btn btn-lg btn-info" type="button" onClick={this.previous}>Previous</button>
          <button className="btn btn-lg btn-primary" type="button" onClick={this.next}>Next</button>
        </div>
      </div>
    );
  }
});

export default DatePicker;

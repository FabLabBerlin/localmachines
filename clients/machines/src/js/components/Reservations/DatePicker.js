var getters = require('../../getters');
var moment = require('moment');
var React = require('react');
var reactor = require('../../reactor');
var ReservationsActions = require('../../actions/ReservationsActions');


class Month {
  // monthId from 0 to 11
  constructor(monthId, year) {
    this._monthId = monthId;
    this._year = year;
    this._firstDay = moment([year, monthId, 1]);
    if (!this._firstDay.isValid()) {
      console.log('this._firstDay is invalid!!!');
    }
    this._lastDay = this._firstDay.clone()
                                  .add(1, 'month')
                                  .subtract(1, 'day');
  }

  static getCurrentMonth() {
    return new Month(moment().month(), moment().year());
  }

  getNextMonth() {
    var t = this.lastDay().add(1, 'day');
    return new Month(t.month(), t.year());
  }

  firstDay() {
    return this._firstDay.clone();
  }

  lastDay() {
    return this._lastDay.clone();
  }

  toString() {
    return this.firstDay().format('MMMM');
  }

  weeks() {
    return this.lastDay().week() - this.firstDay().week() + 1;
  }
}


var DayView = React.createClass({
  handleClick() {
    if (!this.props.header && !this.props.empty && !this.props.notAvailable) {
      var date = this.props.moment;
      console.log('handleClick: date=', date);
      ReservationsActions.createSetDate({ date });
    }
  },

  render() {
    if (this.props.header) {
      return (
        <div className="date-picker-heading">
          {this.props.day}
        </div>
      );
    } else if (this.props.empty) {
      return <div className="date-picker-day date-picker-empty"/>;
    } else {
      var className = 'date-picker-day';
      if (this.props.notAvailable) {
        className += ' date-picker-not-available';
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
  render() {
    var selectedDate = reactor.evaluateToJS(getters.getNewReservation).date;
    var k = 0;
    var month = this.props.month;
    var days = [];
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
                 notAvailable={day.isBefore(moment()) || day.day() === 0}
        />
      );
    }
    while (_.last(weeks).length < 7) {
      _.last(weeks).push(<DayView key={k++} empty={true}/>);
    }

    return (
      <div className="date-picker">
        <h4 className="h4">{month.toString()}</h4>
        <div className="date-picker-week">
          {_.map(['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'], function(heading, i) {
            return <DayView key={i} day={heading} header={true}/>;
          })}
        </div>
        {_.map(weeks, function(w, i) {
          return (
            <div key={i} className="date-picker-week">
              {w}
            </div>
          );
        })}
        <div className="date-picker-finish"/>
      </div>
    );
  }
});


var DatePicker = React.createClass({
  previous() {
    ReservationsActions.previousStep();
  },

  next() {
    ReservationsActions.nextStep();
  },

  render() {
    var currentMonth = Month.getCurrentMonth();
    var nextMonth = currentMonth.getNextMonth();
    return (
      <div className={this.props.className}>
        <h3 className="h3">Select Date</h3>
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

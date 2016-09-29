var _ = require('lodash');
var moment = require('moment');
var React = require('react');


const MONDAY = 1;
const TUESDAY = 2;
const WEDNESDAY = 3;
const THURSDAY = 4;
const FRIDAY = 5;
const SATURDAY = 6;
const SUNDAY = 7;

var DayHeader = React.createClass({
  render() {
    const isCurrentDay = (this.props.day.format('YYYY-MM-MM') ===
      moment().format('YYYY-MM-MM')) ?
      'r-day-current' : '';

    return (
      <div className="r-day-header row">
        <div className="r-day-header-weekday col-xs-3">
          {this.props.day.format('ddd')}
        </div>
        <div className={'r-day-header-number  col-xs-6 ' + isCurrentDay}>
          {this.props.day.format('D')}
        </div>
        <div className="r-day-header-weekday col-xs-3"/>
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


var Day = React.createClass({
  render() {


    return (
      <div className="r-day">
        <DayHeader day={this.props.day}/>
        {_.map(Array(20), () => {
          return <Slot/>;
        })}
      </div>
    );
  }
});


var Week = React.createClass({
  render() {
    return (
      <div className="r-week">
        <Day day={this.props.startDay.clone()}/>
        <Day day={this.props.startDay.clone().add(1, 'day')}/>
        <Day day={this.props.startDay.clone().add(2, 'day')}/>
        <Day day={this.props.startDay.clone().add(3, 'day')}/>
        <Day day={this.props.startDay.clone().add(4, 'day')}/>
        <Day day={this.props.startDay.clone().add(5, 'day')}/>
        <Day day={this.props.startDay.clone().add(6, 'day')}/>
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
        <Week startDay={startDay}/>
      </div>
    );
  }
});

export default Calendar;

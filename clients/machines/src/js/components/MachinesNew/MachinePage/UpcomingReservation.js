var moment = require('moment');
var React = require('react');


function toHHMMSS(d) {
  if (d < 0) {
    return undefined;
  }
  var h = Math.floor(d / 3600);
  var m = Math.floor(d % 3600 / 60);
  var s = Math.floor(d % 3600 % 60);

  return ((h > 0 ? h + ':' + (m < 10 ? '0' : '') : '') + m + ':' + (s < 10 ? '0' : '') + s);
}

var UpcomingReservation = React.createClass({
  render() {
    if (this.props.upcomingReservation) {
      console.log('this.props.upcomingReservation=', this.props.upcomingReservation.toJS());
      const timeStart = moment(this.props.upcomingReservation.get('TimeStart')).unix();
      const delta = timeStart - moment().unix();

      return (
        <div id="m-upcoming-reservation">
          <div id="m-upcoming-reservation-icon"/>
          <div>Next</div>
          <div>Reservation:</div>
          <div className="m-timer">
            {toHHMMSS(delta)}
          </div>
        </div>
      );
    } else {
      return <div/>;
    }
  }
});

export default UpcomingReservation;

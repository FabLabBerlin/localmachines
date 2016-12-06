var React = require('react');
var ReservationTimer = require('../ReservationTimer');


var UpcomingReservation = React.createClass({
  render() {
    if (this.props.upcomingReservation) {
      return (
        <div id="m-upcoming-reservation" className="m-indicator">
          <div id="m-upcoming-reservation-icon"/>
          <div>Next</div>
          <div>Reservation:</div>
          <ReservationTimer reservation={this.props.upcomingReservation.toJS()}/>
        </div>
      );
    } else {
      return <div/>;
    }
  }
});

export default UpcomingReservation;

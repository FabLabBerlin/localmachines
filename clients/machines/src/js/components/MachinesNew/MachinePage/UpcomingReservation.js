var React = require('react');
var ReservationTimer = require('../../MachinePage/Machine/ReservationTimer');


var UpcomingReservation = React.createClass({
  render() {
    if (this.props.upcomingReservation) {
      return (
        <div id="m-upcoming-reservation">
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

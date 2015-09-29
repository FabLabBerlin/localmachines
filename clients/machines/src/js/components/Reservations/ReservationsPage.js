var getters = require('../../getters');
var Navigation = require('react-router').Navigation;
var NewReservation = require('./NewReservation');
var Nuclear = require('nuclear-js');
var React = require('react');
var reactor = require('../../reactor');
var ReservationActions = require('../../actions/ReservationsActions');
var toImmutable = Nuclear.toImmutable;


var ReservationsPage = React.createClass({
  mixins: [ Navigation, reactor.ReactMixin ],

  statics: {
    willTransitionTo(transition) {
      const isLogged = reactor.evaluateToJS(getters.getIsLogged);
      if(!isLogged) {
        transition.redirect('login');
      }
    }
  },

  clickCreate() {
    ReservationActions.createEmpty();
  },

  getDataBindings() {
    return {
      newReservation: getters.getNewReservation,
      reservations: getters.getReservations
    };
  },

  render() {
    if (this.state.newReservation) {
      return <NewReservation/>;
    } else {
      return (
        <div className="container">
          <h3>Reservations</h3>
          <table>
            <thead>
              <th>Name</th>
            </thead>
            <tbody>
              {_.each(this.state.reservations, (reservation) => {
                return (
                  <tr>
                    <td>Reservation</td>
                  </tr>
                );
              })}
            </tbody>
          </table>
          <button className="btn btn-lg btn-primary" onClick={this.clickCreate}>
            Create
          </button>
        </div>
      );
    }
  }
});

export default ReservationsPage;

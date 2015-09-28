var actionTypes = require('../../actionTypes');
var getters = require('../../getters');
var Navigation = require('react-router').Navigation;
var Nuclear = require('nuclear-js');
var React = require('react');
var reactor = require('../../reactor');
var toImmutable = Nuclear.toImmutable;


var ReservationsPage = React.createClass({
  mixins: [ Navigation ],

  statics: {
    willTransitionTo(transition) {
      const isLogged = reactor.evaluateToJS(getters.getIsLogged);
      if(!isLogged) {
        transition.redirect('login');
      }
    }
  },

  render() {
    return (
      <div className="container">
        <h3>Reservations</h3>
        <button className="btn btn-lg btn-primary">
          Create
        </button>
      </div>
    );
  }
});

export default ReservationsPage;

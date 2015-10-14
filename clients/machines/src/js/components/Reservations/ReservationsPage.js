var getters = require('../../getters');
var MachineActions = require('../../actions/MachineActions');
var moment = require('moment');
var Navigation = require('react-router').Navigation;
var NewReservation = require('./NewReservation');
var Nuclear = require('nuclear-js');
var React = require('react');
var reactor = require('../../reactor');
var ReservationRulesActions = require('../../actions/ReservationRulesActions');
var ReservationsActions = require('../../actions/ReservationsActions');
var toImmutable = Nuclear.toImmutable;


function formatDate(date) {
  date = moment(date);
  return date.format('DD. MMM YYYY');
}

function formatTime(date) {
  date = moment(date);
  return date.format('HH:mm');
}

var ReservationsTable = React.createClass({
  mixins: [ reactor.ReactMixin ],

  componentWillMount() {
    const uid = reactor.evaluateToJS(getters.getUid);
    MachineActions.apiGetUserMachines(uid);
    ReservationsActions.load();
    ReservationRulesActions.load();
  },

  getDataBindings() {
    return {
      machinesById: getters.getMachinesById,
      reservations: getters.getReservations
    };
  },

  render() {
    const uid = reactor.evaluateToJS(getters.getUid);
    if (this.state.reservations && this.state.machinesById) {
      return (
        <div className="table-responsive">
          <table className="table table-stripped table-hover">
            <thead>
              <th>Machine</th>
              <th>Date</th>
              <th>From</th>
              <th>To</th>
              <th>Created</th>
            </thead>
            <tbody>
              {_.map(this.state.reservations.toArray(), (reservation, i) => {
                const machineId = reservation.get('MachineId');
                const machine = this.state.machinesById.get(machineId);
                if (machine && reservation.get('UserId') === uid) {
                  return (
                    <tr key={i}>
                      <td>{machine.get('Name')}</td>
                      <td>{formatDate(reservation.get('TimeStart'))}</td>
                      <td>{formatTime(reservation.get('TimeStart'))}</td>
                      <td>{formatTime(reservation.get('TimeEnd'))}</td>
                      <td>{formatDate(reservation.get('Created'))}</td>
                    </tr>
                  );
                } else {
                  console.log('no machine for id ', machineId);
                  console.log('machinesById:', this.state.machinesById);
                }
              })}
            </tbody>
          </table>
        </div>
      );
    } else {
      return <div>'Loading reservations...'</div>;
    }
  }
});


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
    ReservationsActions.createEmpty();
  },

  getDataBindings() {
    return {
      newReservation: getters.getNewReservation
    };
  },

  render() {
    if (this.state.newReservation) {
      return <NewReservation/>;
    } else {
      return (
        <div className="container">
          <h3>My Reservations</h3>
          <ReservationsTable/>
          <hr/>
          <div className="pull-right">
            <button className="btn btn-lg btn-primary" onClick={this.clickCreate}>
              Create
            </button>
          </div>
        </div>
      );
    }
  }
});

export default ReservationsPage;

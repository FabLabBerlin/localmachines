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

// https://github.com/HubSpot/vex/issues/72
var vex = require('vex-js'),
VexDialog = require('vex-js/js/vex.dialog.js');
vex.defaultOptions.className = 'vex-theme-custom';


function formatDate(date) {
  date = moment(date);
  return date.format('DD MMM YYYY');
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

  deleteReservation(reservationId) {
    VexDialog.buttons.YES.text = 'Yes';
    VexDialog.buttons.NO.text = 'No';

    VexDialog.confirm({
      message: 'Do you really want to delete this reservation?',
      callback: function(confirmed) {
        if (confirmed) {
          console.log(reservationId);
          ReservationsActions.deleteReservation(reservationId);
        }
        $('.vex').remove();
        $('body').removeClass('vex-open');
      }.bind(this)
    });
  },

  render() {
    console.log('Reservations page render');
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
              <th>
                <div className="pull-right">
                  Options
                </div>
              </th>
            </thead>
            <tbody>
              {_.map(this.state.reservations.toArray(), (reservation, i) => {
                const machineId = reservation.get('MachineId');
                const machine = this.state.machinesById.get(machineId);
                const reservationId = reservation.get('Id');
                if (machine && reservation.get('UserId') === uid) {
                  return (
                    <tr key={i}>
                      <td>{machine.get('Name')}</td>
                      <td>{formatDate(reservation.get('TimeStart'))}</td>
                      <td>{formatTime(reservation.get('TimeStart'))}</td>
                      <td>{formatTime(reservation.get('TimeEnd'))}</td>
                      <td>{formatDate(reservation.get('Created'))}</td>
                      <td>
                        <button
                          type="button"
                          className="btn btn-danger btn-ico pull-right"
                          onClick={this.deleteReservation.bind(this, reservationId)}>
                          <i className="fa fa-remove"></i>
                        </button>
                      </td>
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

  componentWillUnmount() {
    ReservationsActions.createDone();
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
            <button 
              className="btn btn-lg btn-primary" 
              onClick={this.clickCreate}>
              Create
            </button>
          </div>
        </div>
      );
    }
  }
});

export default ReservationsPage;

var _ = require('lodash');
var $ = require('jquery');
var getters = require('../../getters');
var LoaderLocal = require('../LoaderLocal');
var LocationActions = require('../../actions/LocationActions');
var LocationGetters = require('../../modules/Location/getters');
var MachineActions = require('../../actions/MachineActions');
var Machines = require('../../modules/Machines');
var moment = require('moment');
var Navigation = require('react-router').Navigation;
var NewReservation = require('./NewReservation');
var Nuclear = require('nuclear-js');
var React = require('react');
var reactor = require('../../reactor');
var ReservationRulesActions = require('../../actions/ReservationRulesActions');
var ReservationActions = require('../../actions/ReservationActions');
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


var TableRow = React.createClass({
  cancelReservation(reservationId) {
    VexDialog.buttons.YES.text = 'Yes';
    VexDialog.buttons.NO.text = 'No';

    VexDialog.confirm({
      message: 'Do you really want to cancel this reservation?',
      callback(confirmed) {
        if (confirmed) {
          ReservationActions.cancelReservation(reservationId);
        }
        $('.vex').remove();
        $('body').removeClass('vex-open');
      }
    });
  },

  render() {
    const reservation = this.props.reservation;
    const machine = this.props.machine;
    const reservationId = reservation.get('Id');

    const reservationStart = moment( reservation.get('TimeStart') );
    const reservationEnd = moment( reservation.get('TimeEnd') );
    const now = moment();

    const isPast = reservationEnd.isBefore(now);
    const isToday = (reservationStart.date() === now.date()) &&
      (reservationStart.year() === now.year()) &&
      (reservationStart.month() === now.month());

    const rowClassName = (isPast || reservation.get('Cancelled')) ? 
      'reservation disabled' : 
      'reservation'; 


    return (
      <tr className={rowClassName}>
        <td>{machine.get('Name')}</td>
        <td>{formatDate(reservation.get('TimeStart'))}</td>
        <td>{formatTime(reservation.get('TimeStart'))}</td>
        <td>{formatTime(reservation.get('TimeEnd'))}</td>
        <td>{formatDate(reservation.get('Created'))}</td>
        <td>
          {(!isPast && !isToday) ?
            (!reservation.get('Cancelled') ?
              <button
                type="button"
                className="btn btn-danger btn-ico pull-right"
                onClick={this.cancelReservation.bind(this, reservationId)}>
                <i className="fa fa-remove"></i>
              </button>
            : 'Cancelled')
          : ''}
        </td>
      </tr>
    );
  }
});


var ReservationsTable = React.createClass({
  mixins: [ reactor.ReactMixin ],

  componentWillMount() {
    const locationId = reactor.evaluateToJS(LocationGetters.getLocationId);
    const uid = reactor.evaluateToJS(getters.getUid);
    MachineActions.apiGetUserMachines(locationId, uid);
    ReservationActions.load();
    ReservationRulesActions.load(locationId);
    LocationActions.loadUserLocations(uid);
  },

  getDataBindings() {
    return {
      machinesById: Machines.getters.getMachinesById,
      reservations: getters.getReservations
    };
  },

  render() {
    const uid = reactor.evaluateToJS(getters.getUid);
    if (this.state.reservations && this.state.machinesById) {
      if (this.state.machinesById.size === 0) {
        return <div/>;
      }
      return (
        <div className="table-responsive">
          <table className="table table-stripped table-hover">
            <thead>
              <tr>
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
              </tr>
            </thead>
            <tbody>
              {_.map(this.state.reservations.toArray(), (reservation, i) => {
                const machineId = reservation.get('MachineId');
                const machine = this.state.machinesById.get(machineId);

                if (machine && reservation.get('UserId') === uid) {
                  return <TableRow key={i}
                                   machine={machine}
                                   reservation={reservation}/>;
                }
              })}
            </tbody>
          </table>
        </div>
      );
    } else {
      return <LoaderLocal/>;
    }
  }
});


var ReservationsPage = React.createClass({
  mixins: [ Navigation, reactor.ReactMixin ],

  componentWillUnmount() {
    ReservationActions.newReservation.done();
  },

  clickCreate() {
    ReservationActions.newReservation.create();
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

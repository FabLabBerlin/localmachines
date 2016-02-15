var _ = require('lodash');
var DatePicker = require('./DatePicker');
var getters = require('../../getters');
var MachineActions = require('../../actions/MachineActions');
var moment = require('moment');
var React = require('react');
var reactor = require('../../reactor');
var ReservationsActions = require('../../actions/ReservationsActions');
var TimePicker = require('./TimePicker');
var UserActions = require('../../actions/UserActions');
var toastr = require('../../toastr');


var MachinePricing = React.createClass({

  render() {
    var hourlyPrice = this.props.machine.get('ReservationPriceHourly');
    if (_.isNumber(hourlyPrice)) {
      hourlyPrice = (
        <p><b>Price:</b> €{(hourlyPrice / 2).toFixed(2)} per 30 minutes</p>
      );
    }
    return (
      <div>
        <div className="reservations-machine-price">{hourlyPrice}</div>
        <div>
          The reservation price is on top of the Machine Time.
        </div>
      </div>
    );
  }

});


var SelectMachine = React.createClass({

  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      machines: getters.getMachines,
      machinesById: getters.getMachinesById,
      newReservation: getters.getNewReservation
    };
  },

  handleChange() {
    this.setMachine();
  },

  render() {
    if (this.state.machines.length !== 0) {
      var selectedMachineId;
      var machinePricing;
      if (this.state.newReservation.get('machineId')) {
        selectedMachineId = this.state.newReservation.get('machineId');
        machinePricing = <MachinePricing machine={this.state.machinesById.get(selectedMachineId)}/>;
      }

      return (
        <div className={this.props.className}>
          <h3>Select machine</h3>
          <div>
            <select 
              className="form-control" 
              ref="selection" 
              onChange={this.handleChange} 
              value={selectedMachineId}>
              
              <option value="0">Please select a machine</option>
              
              {_.map(this.state.machines.toArray(), function(machine){
                if (_.isNumber(machine.get('ReservationPriceHourly'))) {
                  return (
                    <option value={machine.get('Id')}>
                      {machine.get('Name')}
                    </option>
                  );
                }
              })}
            </select>
          </div>
          {machinePricing}
          <hr/>
          <div className="pull-right">
            <button 
              className="btn btn-lg btn-info" 
              type="button" 
              onClick={this.cancel}>Cancel</button>

            <button 
              className="btn btn-lg btn-primary" 
              type="button" 
              onClick={this.next}>Next</button>
          </div>
        </div>
      );
    } else {
      return (
        <div>Loading machines...</div>
      );
    }
  },

  cancel() {
    ReservationsActions.newReservation.done();
  },

  next() {
    this.setMachine();
    if (this.state.newReservation.get('machineId')) {
      ReservationsActions.newReservation.nextStep();
    }
  },

  setMachine() {
    var mid = this.refs.selection.getDOMNode().value;
    if (mid) {
      mid = parseInt(mid);
      if (mid) {
        ReservationsActions.newReservation.setMachine({ mid });
      } else {
        toastr.error('No machine selected');
      }
    }
  }
});


var SuccessMsg = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      machinesById: getters.getMachinesById,
      newReservation: getters.getNewReservation,
      newReservationPrice: getters.getNewReservationPrice,
      from: getters.getNewReservationFrom,
      to: getters.getNewReservationTo
    };
  },

  handleClick() {
    ReservationsActions.newReservation.done();
  },

  render() {
    var newReservation = this.state.newReservation;
    const machineId = newReservation.get('machineId');
    const machine = this.state.machinesById.get(machineId);
    const date = moment(this.state.from).format('DD MMM YYYY');
    const timeFrom = moment(this.state.from).format('HH:mm');
    const timeTo = moment(this.state.to).format('HH:mm');
    var containerClassName = 'reservation-confirmed ' + this.props.className;

    return (
      <div className={containerClassName}>
        <h3><i className="fa fa-check-circle-o"></i> Reservation confirmed</h3>

        <div>
          <p><b>Machine:</b> {machine && machine.get('Name')}</p>
          <p><b>Date:</b> {date}</p>
          <p><b>Time:</b> {timeFrom}—{timeTo}</p>
          <p><b>Total price:</b> €{(this.state.newReservationPrice || 0).toFixed(2)}</p>
        </div>
        <div>
          The reservation price is on top of the Machine Time.
        </div>

        <hr/>
        <div className="pull-right">
          <button className="btn btn-lg btn-primary" type="button" onClick={this.handleClick}>
            Continue
          </button>
        </div>
      </div>
    );
  }
});


var NewReservation = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      machines: getters.getMachines,
      machinesById: getters.getMachinesById,
      newReservation: getters.getNewReservation
    };
  },

  componentWillMount() {
    const locationId = reactor.evaluateToJS(getters.getLocation).Id;
    const uid = reactor.evaluateToJS(getters.getUid);
    UserActions.fetchUser(uid);
    MachineActions.apiGetUserMachines(locationId, uid);
  },

  render() {
    if (this.state.machines && this.state.machinesById) {
      var dialog;
      var newReservation = this.state.newReservation;
      var selectedTimes = 0;
      if (newReservation.get('times')) {
        this.state.newReservation.get('times').forEach(function(t) {
          if (t.get('selected')) {
            selectedTimes++;
          }
        });
      }
      switch (newReservation.get('step')) {
      case ReservationsActions.STEP_SET_MACHINE:
        dialog = <SelectMachine className="reservations-new-dialog"/>;
        break;
      case ReservationsActions.STEP_SET_DATE:
        dialog = <DatePicker className="reservations-new-dialog"/>;
        break;
      case ReservationsActions.STEP_SET_TIME:
        dialog = <TimePicker className="reservations-new-dialog"/>;
        break;
      case ReservationsActions.STEP_SUCCESS:
        dialog = <SuccessMsg className="reservations-new-dialog"/>;
        break;
      }
      return (
        <div id="reservations-new" className="container">
          {dialog}
        </div>
      );
    }
  }
});

export default NewReservation;

var DatePicker = require('./DatePicker');
var getters = require('../../getters');
var MachineActions = require('../../actions/MachineActions');
var moment = require('moment');
var React = require('react');
var reactor = require('../../reactor');
var ReservationsActions = require('../../actions/ReservationsActions');
var TimePicker = require('./TimePicker');
var UserActions = require('../../actions/UserActions');


var SelectMachine = React.createClass({

  mixins: [ reactor.ReactMixin ],

  componentWillMount() {
    const uid = reactor.evaluateToJS(getters.getUid);
    UserActions.getUserInfoFromServer(uid);
    MachineActions.apiGetUserMachines(uid);
  },

  getDataBindings() {
    return {
      machineInfo: getters.getMachineInfo
    };
  },

  render() {
    if (this.state.machineInfo.length !== 0) {
      return (
        <div className={this.props.className}>
          <h3 className="h3">Select Machine</h3>
          <div>
            <select className="form-control" ref="selection">
              {_.map(this.state.machineInfo.toArray(), function(machine){
                return (
                  <option value={machine.get('Id')}>
                    {machine.get('Name')}
                  </option>
                );
              })}
            </select>
          </div>
          <hr/>
          <div className="pull-right">
            <button className="btn btn-lg btn-info" type="button" onClick={this.cancel}>Cancel</button>
            <button className="btn btn-lg btn-primary" type="button" onClick={this.setMachine}>Next</button>
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
    ReservationsActions.createDone();
  },

  setMachine() {
    var mid = this.refs.selection.getDOMNode().value;
    mid = parseInt(mid);
    ReservationsActions.createSetMachine({ mid });
  }
});


var SuccessMsg = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      newReservation: getters.getNewReservation,
      from: getters.getNewReservationFrom,
      to: getters.getNewReservationTo
    };
  },

  handleClick() {
    ReservationsActions.createDone();
  },

  render() {
    var newReservation = this.state.newReservation;
    const date = moment(this.state.from).format('DD. MMM YYYY');
    const from = moment(this.state.from).format('HH:mm');
    const to = moment(this.state.to).format('HH:mm');
    return (
      <div className={this.props.className}>
        <h3 className="h3">Your booking is confirmed.</h3>
        <h4>Time:</h4>
        <div>
          {date}
        </div>
        <div>
          {from} - {to}
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
      newReservation: getters.getNewReservation
    };
  },

  render() {
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
});

export default NewReservation;

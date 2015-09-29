var $ = require('jquery');
var getters = require('../../getters');
var MachineActions = require('../../actions/MachineActions');
var React = require('react');
var reactor = require('../../reactor');
var ReservationsActions = require('../../actions/ReservationsActions');
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
    }
  },

  render() {
    if (this.state.machineInfo.length !== 0) {
      return (
        <div>
          <div>Select Machine</div>
          <div>
            <select ref="selection">
              {_.map(this.state.machineInfo.toArray(), function(machine){
                return (
                  <option value={machine.get('Id')}>
                    {machine.get('Name')}
                  </option>
                );
              })}
            </select>
          </div>
          <button type="button" onClick={this.setMachine}>Next</button>
        </div>
      );
    } else {
      return (
        <div>Loading machines...</div>
      );
    }
  },

  setMachine() {
    var mid = this.refs.selection.getDOMNode().value;
    mid = parseInt(mid);
    ReservationsActions.createSetMachine({ mid });
  }
});


var SelectDate = React.createClass({
  render() {
    return (
      <div>
        <div>Select Date</div>
        <input type="text" placeholder="YYYY-MM-DD" ref="date"/>
        <button type="button" onClick={this.setDate}>Next</button>
      </div>
    );
  },

  setDate() {
    var date = this.refs.date.getDOMNode().value;
    ReservationsActions.createSetDate({ date });
  }
});


var SelectTimeRange = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      times: getters.getNewReservationTimes
    };
  },

  render() {
    return (
      <div>
        <div>Select Times</div>
        <div ref="times">
          {_.map(this.state.times.toJS(), (t, i) => {
            return (
              <div key={i}>
                <label>
                  <input
                    type="checkbox"
                  />
                  {t.start.format('HH:mm')} - {t.end.format('HH:mm')}
                </label>
              </div>
            );
          })}
        </div>
        <button type="button" onClick={this.setTimes}>Next</button>
      </div>
    );
  },

  setTimes() {
    var times = this.state.times.toJS();
    $(this.refs.times.getDOMNode()).find('input').each(function(i, el) {
      times[i].selected = el.checked;
    });
    ReservationsActions.createSetTimes({ times });
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
    if (newReservation && newReservation.get('times')) {
      this.state.newReservation.get('times').forEach(function(t) {
        if (t.get('selected')) {
          selectedTimes++;
        }
      });
    }
    if (!newReservation || !newReservation.get('machineId')) {
      dialog = <SelectMachine/>;
    } else if (!newReservation.get('date')) {
      dialog = <SelectDate/>;
    } else if (selectedTimes === 0) {
      dialog = <SelectTimeRange/>;
    }
    return (
      <div className="container">
        {dialog}
      </div>
    );
  }
});

export default NewReservation;

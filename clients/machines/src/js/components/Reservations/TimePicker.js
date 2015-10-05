var $ = require('jquery');
var getters = require('../../getters');
var React = require('react');
var reactor = require('../../reactor');
var ReservationsActions = require('../../actions/ReservationsActions');


var TimePicker = React.createClass({
  mixins: [ reactor.ReactMixin ],

  getDataBindings() {
    return {
      newReservation: getters.getNewReservation,
      times: getters.getNewReservationTimes
    };
  },

  isRange(times) {
    var lastIsSelected;
    var rangesFound = 0;
    _.each(times, function(t) {
      if (t.selected && !lastIsSelected) {
        rangesFound++;
      }
      lastIsSelected = t.selected;
    });
    return rangesFound === 1;
  },

  render() {
    var machineId = this.state.newReservation.get('machineId');
    return (
      <div className={this.props.className}>
        <h3 className="h3">Select time range</h3>
        <div className="no-select" ref="times">
          {_.map(this.state.times.toJS(), (t, i) => {
            console.log('t[' + i + '] = ', t);
            var className = 'time-picker-time';
            if (!_.includes(t.availableMachineIds, machineId)) {
              className += ' unavailable';
            } else if (t.selected) {
              className += ' selected';
            }
            return (
              <div key={i} className={className}>
                <label>
                  <input
                    checked={t.selected}
                    type="checkbox"
                    onChange={this.setTimes.bind(this, i)}
                  />
                  {t.start.format('HH:mm')} - {t.end.format('HH:mm')}
                </label>
              </div>
            );
          })}
        </div>
        <hr/>
        <div className="pull-right">
          <button className="btn btn-lg btn-info" type="button" onClick={this.previous}>Previous</button>
          <button className="btn btn-lg btn-primary" type="button" onClick={this.submit}>Submit</button>
        </div>
      </div>
    );
  },

  previous() {
    ReservationsActions.previousStep();
  },

  setTimes(lastClickIndex) {
    var times = this.state.times.toJS();
    $(this.refs.times.getDOMNode()).find('input').each(function(i, el) {
      times[i].selected = el.checked;
    });
    if (!this.isRange(times)) {
      _.each(times, function(t, i) {
        t.selected = lastClickIndex === i;
      });
    }
    ReservationsActions.createSetTimes({ times });
  },

  submit() {
    ReservationsActions.createSubmit();
  }
});

export default TimePicker;

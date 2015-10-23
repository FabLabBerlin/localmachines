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
      newReservationPrice: getters.getNewReservationPrice,
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
    var lastIndex = _.reduce(this.state.times.toJS(), (lastIdx, t, i) => {
      return t.selected ? i : lastIdx;
    }, null);

    var containerClassName = 'time-picker ' + this.props.className;

    return (
      <div className={containerClassName}>
        <h3>Select time range</h3>
        <div className="no-select" ref="times">
          {_.map(this.times(), (t, i) => {
            var className = 'time';
            var onChange;
            var pricingInfo;

            if (!_.includes(t.availableMachineIds, machineId)) {
              className += ' unavailable';
            } else {
              onChange = this.setTimes.bind(this, 
                i + this.state.times.length - this.times().length);
              if (t.selected) {
                className += ' selected';
              }
            }

            if ((i + this.state.times.length - this.times().length) === lastIndex) {
              pricingInfo = (
                <div className="total-price">
                  Total price: €{(this.state.newReservationPrice).toFixed(2)}
                </div>
              );
            }

            return (
              <div key={i} className={className}>
                <label>
                  <input
                    checked={t.selected}
                    type="checkbox"
                    onChange={onChange}
                  />
                  <div className="row">
                    <div className="col-md-3">
                    </div>
                    <div className="col-md-6">
                      {t.start.format('HH:mm')} - {t.end.format('HH:mm')}
                    </div>
                    <div className="col-md-3">
                      {pricingInfo}
                    </div>
                  </div>
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

    var first = null;
    var firstIndex;
    var last = null;
    var lastIndex;
    $(this.refs.times.getDOMNode()).find('input').each(function(i, el) {

      // Find first checked element
      if (first === null && el.checked) {
        first = el;
        firstIndex = i;
      } 

      // Find last checked element
      if (el.checked) {
        last = el;
        lastIndex = i;
      }

      times[i + this.state.times.length - this.times().length].selected = false;
    }.bind(this));

    if ((firstIndex && lastIndex) && (firstIndex === lastIndex)) {
      times[firstIndex + this.state.times.length - this.times().length].selected = true;
    } else {
      var doSelect = false;
      $(this.refs.times.getDOMNode()).find('input').each(function(i, el) {
        if (parseInt(i) === parseInt(firstIndex)) {
          doSelect = true;
        }

        // Do the selection between first and last checked time item
        if (doSelect) {
          times[i + this.state.times.length - this.times().length].selected = true;
        } else {
          times[i + this.state.times.length - this.times().length].selected = false;
        }

        if (parseInt(i) === parseInt(lastIndex)) {
          doSelect = false;
        } 
      }.bind(this));
    }

    ReservationsActions.createSetTimes({ times });
  },

  submit() {
    ReservationsActions.createSubmit();
  },

  times() {
    var times = this.state.times.toJS();
    var machineId = this.state.newReservation.get('machineId');
    var anyTimeAvailable = _.reduce(times, (foundAny, t) => {
      return foundAny || _.includes(t.availableMachineIds, machineId);
    }, false);

    if (!anyTimeAvailable) {
      times = [];
      var i = 0;
      var minHour = 10;
      var maxHour = 19;
      if (this.state.newReservation.get('date').isoWeekday() === 6) {
        minHour = 12;
        maxHour = 18;
      }
      for (var tt = this.state.newReservation.get('date').clone().hours(minHour); i < 2 * (maxHour - minHour); tt.add(30, 'm'), i++) {
        times.push({
          start: tt.clone(),
          end: tt.clone().add(30, 'm'),
          selected: false,
          availableMachineIds: []
        });
      }
    } else {
      var someAvailableAlready = false;
      times = _.reduce(times, (newTimes, ttt) => {
        if (!someAvailableAlready && !_.includes(ttt.availableMachineIds, machineId)) {
          return newTimes;
        } else {
          someAvailableAlready = true;
          newTimes.push(ttt);
          return newTimes;
        }
      }, []);
    }
    return times;
  }
});

export default TimePicker;
